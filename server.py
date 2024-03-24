#---------------------------------------------------------------------------------------
# Terminal Server for remotely debugging my Flet apps.
# (C) Motti Bazar
# You are welcome to use this freely but if you enhance in any way, please repost back.
# If you distribute the code, please include credit to me.
#---------------------------------------------------------------------------------------

import socket

# Server config
host = 'localhost'       # Insert terminal server IP here
port = 8080

# Create socket object
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

# Bind the socket to a specific address and port
server_socket.bind((host, port))

while True:
    
    server_socket.listen(5)                                   # Listen for connection
    print(f"Server listening on {host}:{port}")
    client_socket, client_address = server_socket.accept()    # Accept connection
    print(f"Connection established with {client_address}")

    while True:
        try:
            data = client_socket.recv(1024)       # Receive data from the client
        except:
            print("Connection probably closed by the client. Exiting...")
            break
        
        if not data:
            print("No data received. Somekind of issue. Exiting...")
            break
            
        print(data.decode())

    # Close the connection
    client_socket.close()
    print("-"*50)

server_socket.close()
