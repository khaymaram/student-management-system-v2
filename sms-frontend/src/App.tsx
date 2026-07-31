import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Toaster } from 'sonner';
import AppLayout from './components/layout/AppLayout';
import { StudentsView } from './components/pages/StudentsView';
import { CoursesView } from './components/pages/CoursesView';
import CourseDetailsView from "./components/pages/CourseDetailsView";
import { ProfessorsView } from './components/pages/ProfessorsView';
import StudentDetailsView from './components/pages/StudentDetailsView';
function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route index element={<Navigate to="/roster" replace />} />
          <Route path="/roster" element={<StudentsView />} />
          <Route path="/roster/:studentId" element={<StudentDetailsView/>}/>
          <Route path="/courses" element={<CoursesView />} />
          <Route path="/courses/:courseCode" element={<CourseDetailsView />}/>
          <Route path="/professors" element={<ProfessorsView />}/>
        </Route>
      </Routes>
      <Toaster
        position="top-right"
        richColors
        toastOptions={{
          style: {
            background: 'var(--bg-elevated)',
            border: '1px solid var(--border)',
            color: 'var(--foreground)',
          },
        }}
      />
    </BrowserRouter>
  );
}

export default App;
