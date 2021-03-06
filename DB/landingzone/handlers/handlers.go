package handlers

import (
	"errors"
	"io"
	"log"
	"net"
	"time"

	"github.com/seadsystem/Backend/DB/landingzone/constants"
)

// Errors
var Timeout = errors.New("Action timed out.")

// Handle a client's request
func HandleRequest(conn net.Conn) {
	log.Println("Got a connection.")

	log.Println("Sending HEAD.")
	conn.Write([]byte(constants.HEAD))

	for {
		log.Println("Reading length header...")
		length_header, err := read_bytes(conn, constants.LENGTH_HEADER_SIZE)
		if err != nil {
			read_error(err)
			break
		}

		data_length := int(length_header[1])

		// Check that we got a length header
		if length_header[0] != 'L' || data_length == 0 {
			log.Println("Invalid length header.")

			// TODO: resync here

			break

		} else {
			log.Printf("Length: %d\n", length_header[1])

			// Get the rest of the packet
			data, err := read_bytes(conn, data_length-constants.LENGTH_HEADER_SIZE)

			if err != nil {
				read_error(err)
				break
			}

			log.Println("Read data:")
			log.Println(string(data))
		}

	}

	conn.Write([]byte("Response"))
	conn.Close()
}

// read_error checks the error and prints an appropriate friendly error message.
func read_error(err error) {
	if err != io.EOF {
		log.Println("Read error:", err)
	} else {
		log.Println("Done reading bytes.")
	}
}

// read_bytes reads the specified number of bytes from the connection with an appropriate time limit.
func read_bytes(conn net.Conn, bytes int) (data []byte, err error) {
	// Setup channels
	data_channel := make(chan []byte, 1)
	error_channel := make(chan error, 1)

	// Initiate read in new go routine
	go func() {
		buffer := make([]byte, bytes)
		n, ierr := conn.Read(buffer)
		if ierr != nil {
			error_channel <- ierr
			return
		}
		if bytes != n {
			error_channel <- io.ErrShortWrite
			return
		}
		data_channel <- buffer
	}()

	// Receive result of read
	select {
	case data = <-data_channel:
		// Read resulted in data
	case err = <-error_channel:
		// Read resulted in an error
	case <-time.After(time.Second * constants.READ_TIME_LIMIT):
		// Read timed out
		err = Timeout
	}

	return
}
