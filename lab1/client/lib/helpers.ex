defmodule Client.Helpers do

  def showCommands() do
    IO.write """
    To subscribe to a topic, enter '%sub% \"topic\">'
    To publish to a specific topic, enter '%pub% \"topic\" \"msg\"'
    To unsubscribe from a topic, enter '%unsub% \"topic\"'
    """
  end

  def parse_cmd(cmd, socket) do
    spl = String.split(cmd)
    json = Poison.encode!(%Message{command: Enum.at(spl, 0), topic: Enum.at(spl, 1), content: Enum.at(spl, 2)})
    "#{json}\n"
  end

  def deserialize(msg) do
    Poison.decode!(msg)
  end

end
