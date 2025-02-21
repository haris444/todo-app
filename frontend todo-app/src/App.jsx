import { useEffect, useState } from "react";

function App() {
    const [todos, setTodos] = useState([]);
    const [newTodo, setNewTodo] = useState("");

    // Fetch todos from backend
    useEffect(() => {
        fetchTodos();
    }, []);

    const fetchTodos = () => {
        fetch("http://localhost:8080/todos")
            .then((res) => res.json())
            .then((data) => setTodos(data))
            .catch((err) => console.error("Error fetching todos:", err));
    };

    // Handle new todo submission
    const addTodo = (e) => {
        e.preventDefault();
        if (!newTodo.trim()) return;

        fetch("http://localhost:8080/todos", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ title: newTodo, done: false }),
        })
            .then((res) => res.json())
            .then((data) => {
                setTodos([...todos, data]); // Update UI
                setNewTodo(""); // Clear input
            })
            .catch((err) => console.error("Error adding todo:", err));
    };

    // Handle checking off a todo
    const toggleTodo = (id, done) => {
        fetch(`http://localhost:8080/todos/${id}`, {
            method: "PUT",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ done: !done }),
        })
            .then((res) => res.json())
            .then((updatedTodo) => {
                setTodos(todos.map(todo =>
                    todo.id === id ? { ...todo, done: updatedTodo.done } : todo
                ));
            })
            .catch((err) => console.error("Error updating todo:", err));
    };

    // Handle deleting a todo
    const deleteTodo = (id) => {
        fetch(`http://localhost:8080/todos/${id}`, {
            method: "DELETE",
        })
            .then(() => {
                setTodos(todos.filter(todo => todo.id !== id)); // Remove from UI
            })
            .catch((err) => console.error("Error deleting todo:", err));
    };

    return (
        <div>
            <h1>To-Do List</h1>

            {/* Add Todo Form */}
            <form onSubmit={addTodo}>
                <input
                    type="text"
                    value={newTodo}
                    onChange={(e) => setNewTodo(e.target.value)}
                    placeholder="New todo..."
                />
                <button type="submit">Add</button>
            </form>

            {/* Todo List */}
            <ul>
                {todos.map((todo) => (
                    <li key={todo.id}>
                        <input
                            type="checkbox"
                            checked={todo.done}
                            onChange={() => toggleTodo(todo.id, todo.done)}
                        />
                        {todo.title} {todo.done ? "✅" : "❌"}
                        <button onClick={() => deleteTodo(todo.id)}>❌Delete</button>
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default App;
