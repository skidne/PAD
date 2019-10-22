defmodule Client do
  def start(port) do
    spawn(fn ->
      case :gen_tcp.connect('localhost', port, [:binary, active: true]) do
        {:ok, socket} ->
          IO.write("Connected to #{port}.\n")
          Client.Helpers.showCommands()
          Client.do_loop(socket)

        {:error, reason} ->
          IO.write("Could not listen: #{reason}\n")
      end
    end)
  end

  def do_loop(socket) do
    packet_to_send = IO.gets('> ') |> Client.Helpers.parse_cmd(socket)
    Client.send(socket, packet_to_send)
    Client.receive(socket)
  end

  def receive(socket) do
    receive do
      {:tcp, ^socket, data} ->
        IO.write("Received packet: #{data}\n")
        Client.do_loop(socket)

      {:tcp_closed, ^socket} ->
        Client.close(socket)
    end
  end

  def send(socket, msg) do
    :gen_tcp.send(socket, msg)
  end

  def close(socket) do
    IO.write("Closing the sender")
    :gen_tcp.close(socket)
    IO.write("Done")
  end
end
