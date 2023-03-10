import json

import requests
from typing import List, Union


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

    def run_command(self, command: str) -> Union[bool, str]:
        response = self.execute_action("eval", f"return {command}")

        if not response:
            return False

        if command in self.command_map and not self.is_undo:
            self.movement_commands.append(command)

        return response

    def execute_action(self, action: str, code: str) -> str:
        r = requests.post(self.command_url, params={"id": self.turtle_id},
                          json={"Action": action, "Code": code})

        try:
            response = r.json()
        except json.JSONDecodeError:
            response = r.text

        return response

    def kill(self) -> str:
        return self.execute_action("kill", "")

    @property
    def command_url(self):
        return f"{self.cnc_server_url}/command"

    def shell(self, command: str) -> str:
        return self.execute_action("shell", command)

    def run_command_n(self, command: str, times: int) -> bool:
        for i in range(times):
            is_success = self.run_command(command)
            if not is_success:
                return False

        return True

    def forward(self, steps: int = 1) -> bool:
        return self.run_command_n("turtle.forward()", steps)

    def dig(self, side: str = Side.right) -> bool:
        return self.run_command(f"turtle.dig('{side}')")

    def back(self, steps: int = 1) -> bool:
        return self.run_command_n("turtle.back()", steps)

    def up(self, steps: int = 1) -> bool:
        return self.run_command_n("turtle.up()", steps)

    def down(self, steps: int = 1) -> bool:
        return self.run_command_n("turtle.down()", steps)

    def get_fuel_level(self) -> bool:
        return self.run_command("turtle.getFuelLevel()")

    def attack(self, side: str = Side.right) -> bool:
        return self.run_command(f"turtle.attack('{side}')")

    def turn_left(self) -> bool:
        return self.run_command("turtle.turnLeft()")

    def turn_right(self) -> bool:
        return self.run_command("turtle.turnRight()")

    @staticmethod
    def command(func):
        def command_wrapper(self):
            return self.run_command(func(self))

        return command_wrapper

    @command
    def inspect(self) -> str:
        return "turtle.inspect()"

    def set_label(self, label: str) -> str:
        return self.run_command(f"os.setComputerLabel('{label}')")
