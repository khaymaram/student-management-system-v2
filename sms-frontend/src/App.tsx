import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { Toaster } from 'sonner';
import AppLayout from './components/layout/AppLayout';
import { StudentsView } from './components/StudentsView';
import { CoursesView } from './components/CoursesView';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route index element={<Navigate to="/roster" replace />} />
          <Route path="/roster" element={<StudentsView />} />
          <Route path="/courses" element={<CoursesView />} />
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
