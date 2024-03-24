import socket
#---------------------------------------------------------------------------------------
# Terminal Client code for remotely debugging my Flet apps.
# (C) Motti Bazar
# You are welcome to use this freely but if you enhance in any way, please repost back.
# If you distribute the code, please include credit to me.
#-----------------------------------------------------------------------------------------------------

# False - print locally using the standard print()
# True  - send the messages to the remote terminal server
send_msg_via_socket = True

terminal_server_ip   = 'localhost'
terminal_server_port = 8080
socket_initialized   = False
client_socket        = None


def showStatusMsg(msg):
    global terminal_server_ip, terminal_server_port, send_msg_via_socket, socket_initialized, client_socket

    if send_msg_via_socket == False:
        print(msg)
        return
    
    # Send via socket
    if socket_initialized == False:
        # Create socket object
        client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        # Connect to the server
        client_socket.connect((terminal_server_ip, terminal_server_port))
        #print(f"Connected to terminal server at {terminal_server_ip}:{terminal_server_port}")
        socket_initialized = True

    client_socket.send(msg.encode())