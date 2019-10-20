defmodule Client do
  require Logger

  def init(port) do
    spawn fn ->
      case :gen_tcp.connect('localhost', port, [:binary, packet: :line, active: false, reuseaddr: true]) do
        {:ok, socket} ->
          Logger.info "Connected to #{port}."
          showCommands()
          Client.do_loop(socket, true)
        {:error, reason} ->
          Logger.error "Could not listen: #{reason}"
      end
    end
  end

  def showCommands() do
    Logger.info "To subscribe to a topic, enter 'sub <topic>'"
    Logger.info "To publish to a specific topic, enter 'pub <topic> <msg>'"
    Logger.info "To unsubscribe from a topic, enter 'unsub <topic>'"
    Logger.info "To close the client, enter 'quit'"
  end

  def do_loop(socket, go) when go == false do
    Logger.info "Client closing..."
    Client.close(socket)
  end

  def do_loop(socket, go) do
    cmd = IO.gets('> ')
    ok = Client.parse_cmd(socket, cmd)
    do_loop(socket, !ok)
  end

  def parse_cmd(socket, cmd) do
    # spl = String.split(cmd)
    Client.send(socket, cmd)
    Client.receive(socket)
    # send formatted to MB
    String.equivalent?(cmd, "quit\n")
  end

  def send(socket, msg) do
    :gen_tcp.send(socket, msg)
  end

  def receive(socket) do
    {:ok, data} = :gen_tcp.recv(socket, 0)
    data
  end

  def close(socket) do
    Logger.info("Closing the sender")
    :gen_tcp.close(socket)
    Logger.info("Done")
  end
end
