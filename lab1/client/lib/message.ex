defmodule Message do
  @derive [Poison.Encoder]
  defstruct [:command, :topic, :content]
end
