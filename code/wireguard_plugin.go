package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

import (
	"github.com/gorilla/mux"
)

var UNIX_PLUGIN_LISTENER = "/state/api/wireguard_plugin"
//var UNIX_PLUGIN_LISTENER = "./http.sock"

type KeyPair struct {
	PrivateKey string
	PublicKey  string
}

type ClientInterface struct {
	PrivateKey string
	Address    string
	DNS        string
}

type ClientPeer struct {
	PublicKey           string
	AllowedIPs          string
	Endpoint            string
	PersistentKeepalive uint
}

type ClientConfig struct {
	Interface ClientInterface
	Peer      ClientPeer
}

// generate a new keypair for a client
func genKeyPair() (KeyPair, error) {
	keypair := KeyPair{}

	cmd := exec.Command("wg", "genkey")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println("wg genkey failed", err)
		return keypair, err
	}

	keypair.PrivateKey = strings.TrimSuffix(string(stdout), "\n")

	cmd = exec.Command("wg", "pubkey")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println("wg pubkey failed", err)
		return keypair, err
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, keypair.PrivateKey+"\n")
	}()

	pubkey, err := cmd.Output()
	if err != nil {
		fmt.Println("wg pubkey failed: bad key", err)
		return keypair, err
	}

	keypair.PublicKey = strings.TrimSuffix(string(pubkey), "\n")

	return keypair, nil
}

// TODO return a new client ip that is not used
func getClientAddress() (string, error) {
	return "192.168.3.4/24", nil
}

func pluginGenKey(w http.ResponseWriter, r *http.Request) {
	keypair, err := genKeyPair()
	if err != nil {
		fmt.Println("wg key failed")
		http.Error(w, "Not found", 404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keypair)
}

// get wireguard endpoint from the environment
func getEndpoint() (string, error) {
	network, ok := os.LookupEnv("WIREGUARD_NETWORK")
	if !ok {
		return "", errors.New("WIREGUARD_NETWORK not set")
	}

	port, ok := os.LookupEnv("WIREGUARD_PORT")
	if !ok {
		return "", errors.New("WIREGUARD_PORT not set")
	}

	endpoint := fmt.Sprintf("%s:%s", network, port)
	return endpoint, nil
}

// get wireguard endpoint from the environment
func getPublicKey() (string, error) {
	pubkey, ok := os.LookupEnv("WIREGUARD_PUBKEY")
	if !ok {
		return "", errors.New("WIREGUARD_PUBKEY not set")
	}

	return pubkey, nil
}

// return config for a client
func pluginGetConfig(w http.ResponseWriter, r *http.Request) {
	config := ClientConfig{}

	keypair, err := genKeyPair()
	if err != nil {
		fmt.Println("wg key failed")
		http.Error(w, "Not found", 404)
		return
	}

	address, err := getClientAddress()
	if err != nil {
		fmt.Println("failed to get client address")
		http.Error(w, "Not found", 404)
		return
	}

	endpoint, err := getEndpoint()
	if err != nil {
		fmt.Println("failed to get endpoint address:", err)
		http.Error(w, "Not found", 404)
		return
	}

	config.Interface.PrivateKey = keypair.PrivateKey
	config.Interface.Address = address
	config.Interface.DNS = "1.1.1.1, 1.0.0.1"

	pubkey, err := getPublicKey()
	if err != nil {
		fmt.Println("failed to get server pubkey:", err)
		http.Error(w, "Not found", 404)
		return
	}

	config.Peer.PublicKey = pubkey
	config.Peer.AllowedIPs = "0.0.0.0/0"
	config.Peer.Endpoint = endpoint
	config.Peer.PersistentKeepalive = 25

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	unix_plugin_router := mux.NewRouter().StrictSlash(true)

	unix_plugin_router.HandleFunc("/genkey", pluginGenKey).Methods("GET")
	unix_plugin_router.HandleFunc("/config", pluginGetConfig).Methods("GET")

	os.Remove(UNIX_PLUGIN_LISTENER)
	unixPluginListener, err := net.Listen("unix", UNIX_PLUGIN_LISTENER)
	if err != nil {
		panic(err)
	}

	pluginServer := http.Server{Handler: logRequest(unix_plugin_router)}

	pluginServer.Serve(unixPluginListener)
}
