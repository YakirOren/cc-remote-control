import time

from ccturtle import Turtle, get_connected_turtles


def main():
    turtles = get_connected_turtles("http://localhost:4000")
    print(turtles)

    t = Turtle("http://localhost:4000", turtles[0])

    t.forward(10)

    t.undo_all()


if __name__ == '__main__':
    main()
