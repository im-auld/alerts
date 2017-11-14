package main

import (
    "fmt"
    "flag"
    "google.golang.org/grpc/metadata"
    "github.com/im-auld/alerts/client"
    "context"
    "os"
    pb "github.com/im-auld/alerts/proto"
)


var sendCommand = flag.NewFlagSet("send", flag.ExitOnError)
var sendUserFlag = sendCommand.Int64("user", -1, "ID of the recipient")
var messageFlag = sendCommand.String("message", "Generic message", "Text message.")
var threadFlag = sendCommand.Int64("thread", -1, "The ID of the thread")
var pathFlag = sendCommand.String("path", "", "The action path")

var seenCommand = flag.NewFlagSet("seen", flag.ExitOnError)
var seenUserFlag = seenCommand.Int64("user", -1, "ID of the recipient")
var uniqFlag = seenCommand.Int64("uniq", -1, "ID of the alert")

var getCommand = flag.NewFlagSet("get", flag.ExitOnError)
var getUserFlag = getCommand.Int64("user", -1, "ID of the recipient")


func newAlertContext(metadataPairs map[string]string) context.Context {
    md := metadata.New(metadataPairs)
    ctx := metadata.NewOutgoingContext(context.Background(), md)
    return ctx
}

func getAlertsForUserContext() context.Context {
    metadataPairs := map[string]string{
        "endpoint": "GetAlertsForUser",
        "caller": "alerts-cli",
    }
    return newAlertContext(metadataPairs)
}

func sendAlertContext() context.Context {
    metadataPairs := map[string]string{
        "endpoint": "SendAlert",
        "caller": "alerts-cli",
    }
    return newAlertContext(metadataPairs)
}

func markAlertSeenContext() context.Context {
    metadataPairs := map[string]string{
        "endpoint": "MarkAlertSeen",
        "caller": "alerts-cli",
    }
    return newAlertContext(metadataPairs)
}

func main() {
    flag.Parse()
    if len(os.Args) == 1 {
        fmt.Println("usage: alerts <command> [<args>]")
        return
    }
    svcHost := "localhost"
    svcPort := "8081"
    c := client.NewAlertClient(svcHost, svcPort)

    switch os.Args[1] {
    case "send":
        sendCommand.Parse(os.Args[2:])
    case "seen":
        seenCommand.Parse(os.Args[2:])
    case "get":
        getCommand.Parse(os.Args[2:])
    default:
        fmt.Printf("%q is not valid command.\n", os.Args[1])
        os.Exit(2)
    }

    if sendCommand.Parsed() {
        alert := &pb.Alert{
            RecipientId: *sendUserFlag,
            ActionPath: *pathFlag,
            ThreadId: *threadFlag,
            Message: *messageFlag,
        }
        resp, err := c.SendAlert(alert, sendAlertContext())
        fmt.Println(resp, err)
    }

    if seenCommand.Parsed() {
        fmt.Println(c.MarkAlertSeen(*seenUserFlag, *uniqFlag, markAlertSeenContext()))
    }

    if getCommand.Parsed() {
        ctx := getAlertsForUserContext()
        fmt.Println(c.GetAlertsForUser(*getUserFlag, ctx))
    }
}
