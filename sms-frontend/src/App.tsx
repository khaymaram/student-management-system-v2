import logo from './assets/logo.png'
import './App.css'
import { useState } from "react";
import "./App.css";
import { AddStudentForm } from "./components/AddStudentForm";
import { StudentsView } from "./components/StudentsView";
import { CoursesView } from './components/CoursesView';

type Tab = "roster" | "courses";

const TABS: { id: Tab; label: string }[] = [
  { id: "roster", label: "Roster" },
  // { id: "add", label: "Add student" },
  { id: "courses", label: "Courses" },
];

function App() {
  // Track which main tab is active so the correct view is shown.
  const [tab, setTab] = useState<Tab>("roster");

  return (
    <>   
    <nav className="tabs">
        {TABS.map((t) => (
          <button
            key={t.id}
            className={`tab${tab === t.id ? " active" : ""}`}
            onClick={() => setTab(t.id)}
          >
            {t.label}
          </button>
        ))}
      </nav>
      <section id="center" className='section'>
        <div>
          <img src={logo} alt="sms logo" className="logo"/>
        </div>
        <div className='title'>
          <h2>Simple Student Management System</h2>
        </div>
      </section>


   

      <section className="tab-panel">
        {tab === "roster" && <StudentsView />}
        {/* {tab === "add" && <AddStudentForm />} */}
        {tab === "courses" && <CoursesView/>}
      </section>

    
    </>
  )
}

export default App
