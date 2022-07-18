package internal

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	pb "chatbox/pkg/chat_services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var messageFormat = "%v >> %v\n %s\n"

type client struct {
	pb.ChatServicesClient
	host        string
	username    string
	currentRoom string
}

func requestContextWithDeadline() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
}

func (c *client) requestRoom(roomName *pb.RoomRequest) {
	ctx, cancel := requestContextWithDeadline()
	defer cancel()
	room, err := c.ChatServicesClient.CreateRoom(ctx, roomName)
	if err != nil {
		log.Printf("client.CreateRoom failed: %v\n", err)
		return
	}
	c.currentRoom = room.RoomId.Id
	fmt.Println(room)
}

func (c *client) viewActiveRooms() {
	ctx, cancel := requestContextWithDeadline()
	defer cancel()
	activeRooms, err := c.ChatServicesClient.GetActiveRooms(ctx, &pb.ActiveRoomsRequest{
		RoomType: 0,
	})
	if err != nil {
		log.Printf("client.GetActiveRooms failed: %v\n", err)
		return
	}
	if c.currentRoom == "" {
		c.currentRoom = activeRooms.Rooms[0].RoomId.Id
	}
	fmt.Println(activeRooms)
}

func (c *client) stream(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := c.ChatServicesClient.RouteChat(ctx)
	if err != nil {
		return err
	}
	defer client.CloseSend()
	log.Println(">>> connected to stream")
	go c.send(client)
	return c.receive(client)
}

func (c *client) send(client pb.ChatServices_RouteChatClient) {
	c.sendJoinBroadcast(client)
	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)
	r, _ := regexp.Compile("^:([a-z0-9-]*)$")
	for {
		select {
		case <-client.Context().Done():
			log.Println("client send loop disconnected")
		default:
			if sc.Scan() {
				if sc.Text() == "" {
					continue
				}
				if r.Match([]byte(sc.Text())) {
					action := r.FindStringSubmatch(sc.Text())
					actRes, err := handleAction(action[1], r)
					if err != nil {
						log.Println(err)
						return
					}
					fmt.Println(actRes)
					return
				}
				if err := client.Send(&pb.ChatMessage{
					DestRoomId: c.currentRoom,
					Message:    sc.Text(),
					StyleType:  0,
					Color:      "",
					Timestamp:  timestamppb.Now(),
					User: &pb.Username{
						Name: c.username,
					},
				}); err != nil {
					log.Printf("%v ** failed to send message: %v\n", time.Now(), err)
					return
				}
			} else {
				log.Printf("%v ** input scanner failure: %v", time.Now(), sc.Err())
			}
		}
	}
}

func handleAction(action string, r *regexp.Regexp) (string, error) {
	formattedAction := r.ReplaceAllFunc([]byte(action), bytes.ToLower)
	switch string(formattedAction) {
	case "help":
		return `
CLIENT COMMANDS
:join - TBD - connect to a chat room, specify the chatroom name.
:view-rooms - TBD - list all active rooms that are joinable including secure rooms.
:leave - TBD - leave the room you joined.
CHAT ROOM COMMANDS
:emotes - list all emotes to send in chat.
		`, nil
	case "emotes":
		return "something", nil
	default:
		return "", fmt.Errorf("action not supported")
	}
}

func (c *client) sendJoinBroadcast(client pb.ChatServices_RouteChatClient) {
	if err := client.Send(&pb.ChatMessage{
		DestRoomId: c.currentRoom,
		Message:    fmt.Sprintf("%s has joined the chat!", c.username),
		StyleType:  0,
		Color:      "",
		Timestamp:  timestamppb.Now(),
		User: &pb.Username{
			Name: c.username,
		},
	}); err != nil {
		log.Println("failed to send joint broadcast to stream.")
		return
	}
}

func (c *client) receive(sc pb.ChatServices_RouteChatClient) error {
	for {
		res, err := sc.Recv()
		if s, ok := status.FromError(err); ok && s.Code() == codes.Canceled {
			log.Println("stream canceled (usually indicates shutdown)")
			return nil
		} else if err == io.EOF {
			log.Println("stream closed by server")
			return nil
		} else if err != nil {
			return err
		} else if res == nil {
			return fmt.Errorf("message is nil")
		} else {
			fmt.Printf(messageFormat, res.User.Name, res.Timestamp.AsTime(), res.Message)
		}
	}
}

func Client(host, username string) *client {
	return &client{
		host:     host,
		username: username,
	}
}

func (c *client) Run(ctx context.Context) error {
	connCtx, cancel := context.WithDeadline(ctx, time.Now().Add(10*time.Second))
	defer cancel()
	conn, err := grpc.DialContext(connCtx, c.host, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("unable to connect to server: %s", string(c.host))
	}
	defer conn.Close()
	c.ChatServicesClient = pb.NewChatServicesClient(conn)
	c.requestRoom(&pb.RoomRequest{
		Name: "test-room",
	})
	c.viewActiveRooms()

	log.Println(">> connected to server: ", c.host)
	err = c.stream(context.Background())
	if err != nil {
		return err
	}
	return nil
}
