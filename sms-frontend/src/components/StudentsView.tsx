// StudentsView.tsx renders the roster UI and handles filtering, editing, and deletion.
import { useState, type FormEvent } from "react";
import { useDeleteStudent, useStudents, useUpdateStudent, useCreateStudent, type StudentFilter } from "../hooks/useStudents";
import { apiErrorMessage } from "../lib/axios";
import type { Student } from "../types";

type FilterMode = "all" | "grade" | "honors" | "search" | "name";

export function StudentsView() {
  // Keep the current filter mode and the text entered by the user in local state.
  const emptyForm = { studentId: "", name: "", grade: "", gpa: "" };

  const [mode, setMode] = useState<FilterMode>("all")
  const [gradeInput, setGradeInput] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [nameInput, setNameInput] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [banner, setBanner] = useState<{ type: "error" | "success"; text: string } | null>(null);

  // const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  // const [addForm, setAddForm] = useState(emptyForm);
  // const [addFieldErrors, setAddFieldErrors] = useState<Partial<Record<keyof StudentInput, string>>>({});

  const filter: StudentFilter =
    mode === "grade" && gradeInput
      ? { type: "grade", grade: Number(gradeInput) }
      : mode === "honors"
        ? { type: "honors" }
        : mode === "search" && searchInput
          ? { type: "search", studentId: Number(searchInput) }
          : mode === "name" && nameInput
            ? { type: "name", name: String(nameInput) }
            : { type: "all" };

  // Fetch students from the backend using the current filter settings.
  const { data: students, isLoading, isError, error } = useStudents(filter);
  const deleteStudent = useDeleteStudent();
  const createStudent = useCreateStudent();

  function handleDelete(student: Student) {
    // Confirm the action before sending a delete request to the API.
    if (!window.confirm(`Remove ${student.name} (ID ${student.studentId}) from the roster`)) return;
    setBanner(null)
    deleteStudent.mutate(student.studentId, {
      onSuccess: () => setBanner({ type: "success", text: `Removed ${student.name} from the roster.` }),
      onError: (err) => setBanner({ type: "error", text: apiErrorMessage(err) }),
    });
  }

  return (
    <div>
      <h2 className="section-title">Student roster</h2>
      <p className="section-sub">View, filter, edit, or remove students on file.</p>

      {banner && <div className={`banner banner-${banner.type}`}>{banner.text}</div>}

      <div className="filter-bar">
        <div className="field">
          <label htmlFor="filter-mode">Show</label>
          <select
            id="filter-mode"
            value={mode}
            onChange={(e) => {
              setMode(e.target.value as FilterMode);
              setBanner(null);
            }}
          >
            <option value="all">All students</option>
            <option value="grade">By grade</option>
            <option value="honors">Honors (GPA ≥ 3.5)</option>
            <option value="search">Search by student ID</option>
            <option value="name">Search by name</option>
          </select>
        </div>

        {mode === "grade" && (
          <div className="field">
            <label htmlFor="grade-filter">Grade</label>
            <input
              id="grade-filter"
              inputMode="numeric"
              placeholder="3"
              value={gradeInput}
              onChange={(e) => setGradeInput(e.target.value)}
            />
          </div>
        )}

        {mode === "search" && (
          <div className="field">
            <label htmlFor="search-filter">Student ID</label>
            <input
              id="search-filter"
              inputMode="numeric"
              placeholder="1001"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
          </div>
        )}
        {mode === "name" && (
          <div className="field">
            <label htmlFor="search-filter">Student Name</label>
            <input
              id="search-filter"
              inputMode="text"
              placeholder="Khaymar Moe"
              value={nameInput}
              onChange={(e) => setNameInput(e.target.value)}
            />
          </div>
        )}
      </div>

      {isLoading && <p className="section-sub">Loading roster…</p>}
      {isError && <div className="banner banner-error">{apiErrorMessage(error)}</div>}

      {!isLoading && !isError && (!students) && (
        <div className="empty-state">
          <span className="stamp">No record found</span>
          <p>No students match this view yet.</p>
        </div>
      )}

      {!isLoading && students && (
        <table className="ledger">
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Grade</th>
              <th>GPA</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {students.map((student) =>
              editingId === student.studentId ? (
                <EditRow
                  key={student.studentId}
                  student={student}
                  onDone={() => setEditingId(null)}
                  onNotify={setBanner}
                />
              ) : (
                <tr key={student.studentId}>
                  <td>
                    <span className="student-id-badge">#{student.studentId}</span>
                  </td>
                  <td>{student.name}</td>
                  <td>{student.grade}</td>
                  <td className={student.gpa >= 3.5 ? "gpa-honors" : undefined}>{student.gpa.toFixed(2)}</td>
                  <td>
                    <div className="row-actions">
                      <button className="btn btn-outline btn-small" onClick={() => setEditingId(student.studentId)}>
                        Edit
                      </button>
                      <button
                        className="btn btn-danger btn-small"
                        onClick={() => handleDelete(student)}
                        disabled={deleteStudent.isPending}
                      >
                        Remove
                      </button>
                    </div>
                  </td>
                </tr>
              ),
            )}
          </tbody>
        </table>
      )}
    </div>
  );
}

function EditRow({
  student,
  onDone,
  onNotify,
}: {
  student: Student;
  onDone: () => void;
  onNotify: (banner: { type: "error" | "success"; text: string }) => void;
}) {
  // Local state lets the edit row behave like a small inline form.
  const [name, setName] = useState(student.name);
  const [grade, setGrade] = useState(String(student.grade));
  const [gpa, setGpa] = useState(String(student.gpa));
  const updateStudent = useUpdateStudent();

  function handleSave(e: FormEvent) {
    // Validate the edited values before sending an update request.
    e.preventDefault();
    const gpaNum = Number(gpa);
    const gradeNum = Number(grade);
    if (!name.trim() || Number.isNaN(gpaNum) || Number.isNaN(gradeNum) || gpaNum < 0 || gpaNum > 4) {
      onNotify({ type: "error", text: "Please enter a valid name, grade, and GPA (0.0–4.0)." });
      return;
    }
    updateStudent.mutate(
      { studentId: student.studentId, input: { studentId: student.studentId, name, grade: gradeNum, gpa: gpaNum } },
      {
        onSuccess: () => {
          onNotify({ type: "success", text: `Updated ${name}'s record.` });
          onDone();
        },
        onError: (err) => onNotify({ type: "error", text: apiErrorMessage(err) }),
      },
    );
  }

  return (
    <tr>
      <td>
        <span className="student-id-badge">#{student.studentId}</span>
      </td>
      <td colSpan={3}>
        <form className="edit-row-form" onSubmit={handleSave}>
          <input name="name" value={name} onChange={(e) => setName(e.target.value)} aria-label="Name" />
          <input name="grade" value={grade} onChange={(e) => setGrade(e.target.value)} aria-label="Grade" />
          <input name="gpa" value={gpa} onChange={(e) => setGpa(e.target.value)} aria-label="GPA" />
          <button type="submit" className="btn btn-primary btn-small" disabled={updateStudent.isPending}>
            Save
          </button>
          <button type="button" className="btn btn-outline btn-small" onClick={onDone}>
            Cancel
          </button>
        </form>
      </td>
      <td></td>
    </tr>
  );
}