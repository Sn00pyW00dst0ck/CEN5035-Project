import { test, expect } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import ServerList from "../MainScreen/ServerList/ServerList.jsx";
import ServerBadge from "../MainScreen/ServerList/ServerBadge/ServerBadge.jsx";

const servers = [
    { id: 1, name: "test1", icon: "public/vite.svg" },
    { id: 2, name: "test2", icon: "public/vite.svg" },
    { id: 3, name: "Test1", icon: "public/vite.svg" },
    { id: 4, name: "Test2", icon: "public/vite.svg" },
    { id: 5, name: "thisIsATest1", icon: "public/vite.svg" },
    { id: 6, name: "Alice", icon: "public/vite.svg" },
];

test('Renders without exceptions', () => {
    expect(() => render(<ServerList servers={servers}/>)).not.toThrow();
});

test("Server list renders all test servers", () => {
    render(<ServerList servers={servers} />);

    // Check filtered results
    expect(screen.queryByText("test1")).toBeInTheDocument();
    expect(screen.queryByText("test2")).toBeInTheDocument();
    expect(screen.queryByText("Test1")).toBeInTheDocument();
    expect(screen.queryByText("Test2")).toBeInTheDocument();
    expect(screen.queryByText("thisIsATest1")).toBeInTheDocument();
    expect(screen.getByText("Alice")).toBeInTheDocument();
});

test("Search filters servers returned", () => {
    render(<ServerList servers={servers} />);

    // Use label "Search" to target the TextField
    const input = screen.getByLabelText("Search");

    fireEvent.input(input, { target: { value: "Alice" } });

    // Check filtered results
    expect(screen.queryByText("test1")).not.toBeInTheDocument();
    expect(screen.queryByText("test2")).not.toBeInTheDocument();
    expect(screen.queryByText("Test1")).not.toBeInTheDocument();
    expect(screen.queryByText("Test2")).not.toBeInTheDocument();
    expect(screen.queryByText("thisIsATest1")).not.toBeInTheDocument();
    expect(screen.getByText("Alice")).toBeInTheDocument();
});
