import requests


class Turtle:
    def __init__(self, cnc_server_url: str, turtle_id: str):
        self.cnc_server_url = cnc_server_url
        self.turtle_id = turtle_id

    def send_command(self, command: str) -> str:
        command_url = f"{self.cnc_server_url}/command"
        r = requests.post(command_url, params={"id": self.turtle_id},
                          json={"Action": "eval", "Code": f"return {command}"})
        return r.text

    def kill(self) -> str:
        command_url = f"{self.cnc_server_url}/command"
        r = requests.post(command_url, params={"id": self.turtle_id}, json={"Action": "kill", "Code": ""})
        return r.text

    def shell(self, command: str) -> str:
        command_url = f"{self.cnc_server_url}/command"
        r = requests.post(command_url, params={"id": self.turtle_id}, json={"Action": "shell", "Code": f"{command}"})
        return r.text

    def forward(self) -> str:
        return self.send_command("turtle.forward()")

    def dig(self, side: str) -> str:
        return self.send_command(f"turtle.dig('{side}')")

    def back(self) -> str:
        return self.send_command("turtle.back()")

    def up(self) -> str:
        return self.send_command("turtle.up()")

    def down(self) -> str:
        return self.send_command("turtle.down()")

    def get_fuel_level(self) -> str:
        return self.send_command("turtle.getFuelLevel()")

    def attack(self, side: str) -> str:
        return self.send_command(f"turtle.attack('{side}')")

    def set_label(self, label: str) -> str:
        self.send_command(f"os.setComputerLabel('{label}')")
