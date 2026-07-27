// AddStudentForm.tsx renders the student creation form.
// It validates user input with Zod and submits new student data via the hook.
import { useState, type FormEvent } from "react";
import { StudentInputSchema, type StudentInput } from "../types";
import { useCreateStudent } from "../hooks/useStudents";
import { apiErrorMessage } from "../lib/axios";

const emptyForm = { studentId: "", name: "", grade: "", gpa: "" };

export function AddStudentForm(){
    // Keep the form values in local state so the UI can update as the user types.
    const [form, setForm] = useState(emptyForm);
    // Store validation messages for each field so the user can fix errors quickly.
    const [fieldErrors, setFieldErrors] = useState<Partial<Record<keyof StudentInput, string>>>({});
    // Show a short success or error message after the create request finishes.
    const [banner, setBanner] = useState<{ type: "error" | "success"; text: string } | null>(null);
    const createStudent = useCreateStudent();

    function handleChange(field: keyof typeof form, value: string) {
        setForm((prev: typeof emptyForm) => ({ ...prev, [field]: value }));
    }

    function handleSubmit(e: FormEvent) {
        e.preventDefault();
        setBanner(null);

        // Validate the form before sending anything to the backend.
        const result = StudentInputSchema.safeParse(form);
        if (!result.success) {
            const errors: Partial<Record<keyof StudentInput, string>> = {};
            for (const issue of result.error.issues) {
                const key = issue.path[0] as keyof StudentInput;
                errors[key] = issue.message;
            }
            setFieldErrors(errors);
            return;
        }
        setFieldErrors({});

        createStudent.mutate(result.data, {
            onSuccess: () => {
                setForm(emptyForm);
                setBanner({ type: "success", text: `Enrolled ${result.data.name} (ID ${result.data.studentId}) in the roster.`});
            },
            onError: (err) => {
                setBanner({ type: "error", text: apiErrorMessage(err)});
            },
        });
    }
    return (
    <div>
      <h2 className="section-title">Add a student</h2>
      <p className="section-sub">Enter a new record for the student roster.</p>

      {banner && <div className={`banner banner-${banner.type}`}>{banner.text}</div>}

      <form className="form-grid" onSubmit={handleSubmit} noValidate>
        <div className="field">
          <label htmlFor="studentId">Student ID</label>
          <input
            id="studentId"
            inputMode="numeric"
            placeholder="1005"
            value={form.studentId}
            onChange={(e) => handleChange("studentId", e.target.value)}
          />
          {fieldErrors.studentId && <span className="field-error">{fieldErrors.studentId}</span>}
        </div>

        <div className="field">
          <label htmlFor="name">Full name</label>
          <input
            id="name"
            placeholder="Khaymar Moe"
            value={form.name}
            onChange={(e) => handleChange("name", e.target.value)}
          />
          {fieldErrors.name && <span className="field-error">{fieldErrors.name}</span>}
        </div>

        <div className="field">
          <label htmlFor="grade">Grade</label>
          <input
            id="grade"
            inputMode="numeric"
            placeholder="3"
            value={form.grade}
            onChange={(e) => handleChange("grade", e.target.value)}
          />
          {fieldErrors.grade && <span className="field-error">{fieldErrors.grade}</span>}
        </div>

        <div className="field">
          <label htmlFor="gpa">GPA</label>
          <input
            id="gpa"
            inputMode="decimal"
            placeholder="3.8"
            value={form.gpa}
            onChange={(e) => handleChange("gpa", e.target.value)}
          />
          {fieldErrors.gpa && <span className="field-error">{fieldErrors.gpa}</span>}
        </div>

        <div className="form-actions">
          <button type="submit" className="btn btn-primary" disabled={createStudent.isPending}>
            {createStudent.isPending ? "Adding…" : "Add student"}
          </button>
        </div>
      </form>
    </div>
  );
}
    

