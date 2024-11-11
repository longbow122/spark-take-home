import React, { useEffect, useState } from 'react';
import './App.css';
import Todo, { TodoType } from './Todo';

function App() {
  const [todos, setTodos] = useState<TodoType[]>([]);

  // Initially fetch todo
  useEffect(() => {
    const fetchTodos = async () => {
      try {
        const todos = await fetch('http://localhost:8080/');
        if (todos.status !== 200) {
          console.log('Error fetching data');
          return;
        }

        setTodos(await todos.json());
      } catch (e) {
        console.log('Could not connect to server. Ensure it is running. ' + e);
      }
    }

    fetchTodos()
  }, []);

  const handleSubmission = async (e: React.FormEvent<HTMLFormElement>) => {
    try {
      const formData = new FormData(e.currentTarget);
      const title = formData.get('title') as string;
      const description = formData.get('description') as string;
      const response = await fetch('http://localhost:8080/', {
        method: 'POST',
        body: JSON.stringify({
          title: title,
          description: description,
        }),
      });
      if (response.status !== 200) {
        console.log('Error submitting data');
        return;
      }

      if (response.ok) {
        const updated = await response.json();
        setTodos([...todos, updated]);
      }
    } catch (e) {
      console.log("Error submitting form data.");
      console.log(e);
    }
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>TODO</h1>
      </header>

      <div className="todo-list">
        {todos.map((todo) =>
          <Todo
            key={todo.title + todo.description}
            title={todo.title}
            description={todo.description}
          />
        )}
      </div>

      <h2>Add a Todo</h2>
      <form onSubmit={handleSubmission}>
        <input id="title" placeholder="Title" name="title" autoFocus={true} required={true}/>
        <input id={"description"} placeholder="Description" name="description" required={true}/>
        <button type={"submit"}>Add Todo</button>
      </form>
    </div>
  );
}

export default App;
