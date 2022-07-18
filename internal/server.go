package internal

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	pb "chatbox/pkg/chat_services"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type chatBoxServer struct {
	pb.UnimplementedChatServicesServer
	activeRooms *pb.RoomList
	roomChat    chan *pb.ChatMessage
	mu          sync.Mutex // protects roomMessageRouters
}

func verifyCtx(ctx context.Context) error {
	if errors.Is(ctx.Err(), context.Canceled) {
		return ctx.Err()
	}
	return nil
}

func (server *chatBoxServer) checkActiveRooms() error {
	if len(server.activeRooms.Rooms) == 0 {
		return fmt.Errorf("no active rooms found")
	}
	return nil
}

func (server *chatBoxServer) findRoom(roomId string, roomName string) (*pb.Room, error) {
	err := server.checkActiveRooms()
	if err != nil {
		return nil, err
	}
	passedId := roomId != ""
	passedName := roomName != ""
	if !passedId && !passedName {
		return nil, fmt.Errorf("no parameters passed")
	}
	if len(server.activeRooms.Rooms) == 0 {
		return nil, fmt.Errorf("no active rooms found")
	}
	for _, room := range server.activeRooms.Rooms {
		if passedId && passedName {
			if room.RoomId.Id == roomId && room.RoomName.Name == roomName {
				return room, nil
			}
		}
		if passedId && room.RoomId.Id == roomId {
			return room, nil
		}
		if passedName && room.RoomName.Name == roomName {
			return room, nil
		}
	}
	if passedName {
		return nil, fmt.Errorf("no room found with name %v", roomName)
	}
	return nil, fmt.Errorf("no room found with id %v", roomId)
}

func (server *chatBoxServer) CreateRoom(ctx context.Context, roomRequest *pb.RoomRequest) (*pb.Room, error) {
	if err := verifyCtx(ctx); err != nil {
		return nil, err
	}
	var activeRooms *pb.RoomList
	var newRoom *pb.Room
	if len(server.activeRooms.Rooms) == 0 {
		newRoom = &pb.Room{
			RoomId: &pb.RoomId{
				Id: uuid.New().String(),
			},
			RoomName:     roomRequest,
			Participants: 0,
			IsSecure:     false,
			RoomStatus:   0,
		}
		activeRooms = &pb.RoomList{
			Rooms: []*pb.Room{newRoom},
		}
		server.activeRooms = activeRooms
	} else {
		for _, room := range server.activeRooms.Rooms {
			if roomRequest.Name == room.RoomName.Name {
				return nil, fmt.Errorf("room name already taken")
			}
		}
		newRoom = &pb.Room{
			RoomId: &pb.RoomId{
				Id: uuid.New().String(),
			},
			RoomName:     roomRequest,
			Participants: 0,
			IsSecure:     false,
			RoomStatus:   0,
		}
		activeRooms = server.activeRooms
		activeRooms.Rooms = append(activeRooms.Rooms, newRoom)
		server.mu.Lock()
		server.activeRooms = activeRooms
		server.mu.Unlock()
	}

	return newRoom, nil
}

func (server *chatBoxServer) GetActiveRooms(ctx context.Context, activeRoomRequest *pb.ActiveRoomsRequest) (*pb.RoomList, error) {
	if err := verifyCtx(ctx); err != nil {
		return nil, err
	}
	lobby := server.activeRooms
	if len(lobby.Rooms) == 0 {
		return nil, fmt.Errorf("the lobby is empty")
	}
	activeRooms := []*pb.Room{}
	for _, room := range lobby.Rooms {
		if room.RoomStatus == activeRoomRequest.RoomType {
			activeRooms = append(activeRooms, room)
		}
	}
	return &pb.RoomList{
			Rooms: activeRooms,
		},
		nil
}

func (server *chatBoxServer) GetRoom(ctx context.Context, roomId *pb.RoomId) (*pb.Room, error) {
	if err := verifyCtx(ctx); err != nil {
		return nil, err
	}
	room, err := server.findRoom(roomId.Id, "")
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (server *chatBoxServer) RouteChat(stream pb.ChatServices_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if in == nil {
			return fmt.Errorf("received message is nil")
		}

		go server.sendMessages(stream)
		room, err := server.findRoom(in.DestRoomId, "")
		if err != nil {
			return err
		}
		if room.RoomStatus != pb.RoomType_ACTIVE {
			return fmt.Errorf("server.RouteChat; stream.Recv: provided destRoomId %v is INACTIVE", in.DestRoomId)
		}
		server.mu.Lock()
		in.Timestamp = timestamppb.Now()
		server.roomChat <- in
		server.mu.Unlock()
	}
}

func (server *chatBoxServer) sendMessages(srv pb.ChatServices_RouteChatServer) {
	for {
		select {
		case <-srv.Context().Done():
			return
		case msg := <-server.roomChat:
			if s, ok := status.FromError(srv.Send(msg)); ok {
				switch s.Code() {
				case codes.OK:
					// noop
					fmt.Println(msg)
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					log.Println("client terminated connection")
				default:
					log.Printf("%v failed to send message to a client: %v", time.Now(), s.Err())
				}
			}
		}
	}
}

func newServer() *chatBoxServer {
	server := &chatBoxServer{
		activeRooms: &pb.RoomList{
			Rooms: []*pb.Room{},
		},
		roomChat: make(chan *pb.ChatMessage),
	}
	return server
}

func Server(port *int) (*grpc.Server, net.Listener) {
	serverAddr := fmt.Sprintf("localhost:%d", *port)
	listen, err := net.Listen("tcp", serverAddr)
	if err != nil {
		log.Fatalln(err)
	}
	var ops []grpc.ServerOption

	grpcServer := grpc.NewServer(ops...)
	pb.RegisterChatServicesServer(grpcServer, newServer())
	log.Printf("I AM LISTENING AT: %d\n", *port)
	return grpcServer, listen
}
