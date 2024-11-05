package main

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "time"

    "google.golang.org/grpc"
    pb "path/to/your/protobuf/package" // Update with your actual protobuf package path
)

// Config represents the configuration structure
type Config struct {
    Token    string `json:"token"`
    Prefix   string `json:"prefix"`
    LightNum string `json:"lightNum"`
    Username string `json:"username"`
    Host     string `json:"host"`
}

// logError logs the error and exits the program if the error is not nil
func logError(err error) {
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
}

// loadConfig loads configuration from config.json
func loadConfig(filename string) (*Config, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("could not open config file: %w", err)
    }
    defer file.Close()

    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("could not read config file: %w", err)
    }

    var config Config
    err = json.Unmarshal(bytes, &config)
    if err != nil {
        return nil, fmt.Errorf("could not unmarshal config data: %w", err)
    }

    return &config, nil
}

// connectToGRPCServer establishes a connection to the gRPC server
func connectToGRPCServer(address string) (*grpc.ClientConn, error) {
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        return nil, fmt.Errorf("did not connect: %w", err)
    }
    return conn, nil
}

// grpcUserToUserInfo performs the first RPC call
func grpcUserToUserInfo(client pb.YourServiceClient, token string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := client.YourRPCMethod(ctx, &pb.YourRequest{Token: token})
    logError(err)
    fmt.Println("Response from RPC1: ", response)
}

// grpcStructToUserAttributes performs the second RPC call
func grpcStructToUserAttributes(client pb.YourServiceClient, token string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    response, err := client.AnotherRPCMethod(ctx, &pb.AnotherRequest{Token: token})
    logError(err)
    fmt.Println("Response from RPC2: ", response)
}

func main() {
    // Load the configuration from config.json
    config, err := loadConfig("config.json")
    logError(err)

    if config.Token == "" {
        log.Fatalf("API token is missing")
    }

    // Set up a connection to the gRPC server
    conn, err := connectToGRPCServer(config.Host + ":50051")
    logError(err)
    defer conn.Close()

    client := pb.NewYourServiceClient(conn)

    // Perform multiple RPC calls
    grpcUserToUserInfo(client, config.Token)
    grpcStructToUserAttributes(client, config.Token)

    // Additional logic if necessary
    fmt.Println("All RPC calls were successful")
}
