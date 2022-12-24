import requests
from typing import List


class Side:
    left = "left"
    right = "right"


def get_connected_turtles(cnc_server_url: str) -> List[str]:
    return requests.get(cnc_server_url + "/sessions").json()


class Turtle:
    def __init__(self, cnc_server_url: str, turtle_id: str):
        self.cnc_server_url = cnc_server_url
        self.turtle_id = turtle_id
        self.movement_commands = []
        self.command_map = {"turtle.forward()": self.back,
                            "turtle.back()": self.forward,
                            "turtle.down()": self.up,
                            "turtle.up()": self.down,
                            "turtle.turnLeft()": self.turn_right,
                            "turtle.turnRight()": self.turn_left
                            }
        self.is_undo = False

    def undo(self):
        command = self.movement_commands.pop()
        self.is_undo = True
        self.command_map[command]()
        self.is_undo = False

    def undo_all(self):
        for i in range(len(self.movement_commands)):
            self.undo()

    def run_command(self, command: str) -> str:
        command_url = f"{self.cnc_server_url}/command"
        r = requests.post(command_url, params={"id": self.turtle_id},
                          json={"Action": "eval", "Code": f"return {command}"})

        if command in self.command_map and r.text == "true" and not self.is_undo:
            self.movement_commands.append(command)

        return r.text

    def kill(self) -> str:
        command_url = f"{self.cnc_server_url}/command"
        r = requests.post(command_url, params={"id": self.turtle_id}, json={"Action": "kill", "Code": ""})
        return r.text

    def shell(self, command: str) -> str:
        command_url = f"{self.cnc_server_url}/command"
        r = requests.post(command_url, params={"id": self.turtle_id}, json={"Action": "shell", "Code": f"{command}"})
        return r.text

    def run_command_n(self, command: str, times: int) -> str:
        for i in range(times):
            output = self.run_command(command)
            if output != "true":
                return output

        return "true"

    def forward(self, steps: int) -> str:
        return self.run_command_n("turtle.forward()", steps)

    def dig(self, side: str) -> str:
        return self.run_command(f"turtle.dig('{side}')")

    def back(self) -> str:
        return self.run_command("turtle.back()")

    def up(self) -> str:
        return self.run_command("turtle.up()")

    def down(self) -> str:
        return self.run_command("turtle.down()")

    def get_fuel_level(self) -> str:
        return self.run_command("turtle.getFuelLevel()")

    def attack(self, side: str) -> str:
        return self.run_command(f"turtle.attack('{side}')")

    def turn_left(self) -> str:
        return self.run_command("turtle.turnLeft()")

    def turn_right(self) -> str:
        return self.run_command("turtle.turnRight()")

    def inspect(self) -> str:
        return self.run_command("turtle.inspect()")

    def set_label(self, label: str) -> str:
        self.run_command(f"os.setComputerLabel('{label}')")
