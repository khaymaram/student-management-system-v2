import { useEffect, useState, type FormEvent } from "react";
import { useDeleteCourse, useCourses, useUpdateCourses, useCreateCourse, type CourseFilter } from "../hooks/useCourses";
import { apiErrorMessage } from "../lib/axios";
import { CourseInputSchema, type Course, type CourseInput } from "../types";

type FilterMode = "all" | "credits" | "title" | "code";

export function CoursesView() {
  // Keep the current filter mode and the text entered by the user in local state.
  const emptyForm = { title: "", credits: "", code: "" };
  const [mode, setMode] = useState<FilterMode>("all")
  const [creditInput, setCreditInput] = useState("");
  const [titleInput, setTitleInput] = useState("");
  const [codeInput, setCodeInput] = useState("");
  const [editingCode, setEditingCode] = useState<string | null>(null);
  const [banner, setBanner] = useState<{ type: "error" | "success"; text: string } | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [addForm, setAddForm] = useState(emptyForm);
  const [addFieldErrors, setAddFieldErrors] = useState<Partial<Record<keyof CourseInput, string>>>({});

  useEffect(() => {
    if (!banner) return;

    const timer = window.setTimeout(() => setBanner(null), 3000);
    return () => window.clearTimeout(timer);
  }, [banner]);

  const filter: CourseFilter =
    mode === "credits" && creditInput
      ? { type: "credits", credits: Number(creditInput) }
      : mode === "title" && titleInput
        ? { type: "title", title: String(titleInput) }
        : mode === "code" && codeInput
          ? { type: "code", code: String(codeInput) }
          : { type: "all" };

  // Fetch courses from the backend using the current filter settings.
  const { data: courses, isLoading, isError, error } = useCourses(filter);
  const deleteCourse = useDeleteCourse();
  const createCourse = useCreateCourse();

  function handleDelete(course: Course) {
    // Confirm the action before sending a delete request to the API.
    if (!window.confirm(`Remove ${course.title} (${course.code}) from the directory`)) return;
    setBanner(null)
    deleteCourse.mutate(course.code, {
      onSuccess: () => setBanner({ type: "success", text: `Removed ${course.title} from the directory.` }),
      onError: (err) => setBanner({ type: "error", text: apiErrorMessage(err) }),
    });
  }

  function handleAddFormChange(field: keyof typeof addForm, value: string) {
    setAddForm((prev) => ({ ...prev, [field]: value }));
  }

  function openAddCourseModal() {
    setAddForm({ ...emptyForm });
    setAddFieldErrors({});
    setBanner(null);
    setIsAddModalOpen(true);
  }

  function addCourse(e: FormEvent) {
    e.preventDefault();
    setBanner(null);

    const result = CourseInputSchema.safeParse(addForm);
    if (!result.success) {
      const errors: Partial<Record<keyof CourseInput, string>> = {};
      for (const issue of result.error.issues) {
        const key = issue.path[0] as keyof CourseInput;
        errors[key] = issue.message;
      }
      setAddFieldErrors(errors);
      return;
    }

    setAddFieldErrors({});
    createCourse.mutate(result.data, {
      onSuccess: () => {
        setAddForm({ ...emptyForm });
        setIsAddModalOpen(false);
        setBanner({ type: "success", text: `Created ${result.data.title} (${result.data.code}) in the directory.` });
      },
      onError: (err) => {
        setBanner({ type: "error", text: apiErrorMessage(err) });
      },
    });
  }


  return (
    <div>
      <h2 className="section-title">Course Directory</h2>
      <p className="section-sub">View, filter, edit, or remove courses.</p>

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
            <option value="all">All courses</option>
            <option value="credits">By credits</option>
            <option value="title">Search by course title</option>
            <option value="code">Search by course code</option>
          </select>
        </div>

        {mode === "credits" && (
          <div className="field">
            <label htmlFor="grade-filter">Credits</label>
            <input
              id="grade-filter"
              inputMode="numeric"
              placeholder="3"
              value={creditInput}
              onChange={(e) => setCreditInput(e.target.value)}
            />
          </div>
        )}

        {mode === "title" && (
          <div className="field">
            <label htmlFor="search-filter">Course Title</label>
            <input
              id="search-filter"
              inputMode="text"
              placeholder="Intro to OOP"
              value={titleInput}
              onChange={(e) => setTitleInput(e.target.value)}
            />
          </div>
        )}
        {mode === "code" && (
          <div className="field">
            <label htmlFor="search-filter">Course Code</label>
            <input
              id="search-filter"
              inputMode="text"
              placeholder="CMSC131"
              value={codeInput}
              onChange={(e) => setCodeInput(e.target.value)}
            />
          </div>
        )}
      </div>

      {isLoading && <p className="section-sub">Loading roster…</p>}
      {isError && <div className="banner banner-error">{apiErrorMessage(error)}</div>}

      {!isLoading && !isError && (!courses) && (
        <div className="empty-state">
          <span className="stamp">No record found</span>
          <p>No courses match this view yet.</p>
        </div>
      )}

      {!isLoading && courses && (
        <table className="ledger">
          <thead>
            <tr>
              <th>Course Code</th>
              <th>Title</th>
              <th>Credits</th>
              
              <th></th>
            </tr>
          </thead>
          <tbody>
            {courses.map((course) =>
              editingCode === course.code ? (
                <EditRow
                  key={course.code}
                  course={course}
                  onDone={() => setEditingCode(null)}
                  onNotify={setBanner}
                />
              ) : (
                <tr key={course.code}>
                  <td><span className="student-id-badge">{course.code}</span></td>
                  <td>{course.title}</td>
                  <td>{course.credits}</td>
                  
                  <td>
                    <div className="row-actions">
                      <button className="btn btn-outline btn-small" onClick={() => setEditingCode(course.code)}>
                        Edit
                      </button>
                      <button
                        className="btn btn-danger btn-small"
                        onClick={() => handleDelete(course)}
                        disabled={deleteCourse.isPending}
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
      <br></br>
      <button type="button" className="btn btn-primary" onClick={openAddCourseModal}>Add a Course</button>

      {isAddModalOpen && (
        <div
          role="dialog"
          aria-modal="true"
          className="modal-backdrop"
          onClick={() => {
            setIsAddModalOpen(false);
            setAddForm({ ...emptyForm });
            setAddFieldErrors({});
          }}
        >
          <div
            className="modal-card"
            onClick={(e) => e.stopPropagation()}
          >
            <div className="modal-header">
              <h3 className="modal-title">Add a Course</h3>
              <button
                type="button"
                className="btn btn-outline btn-small"
                onClick={() => {
                  setIsAddModalOpen(false);
                  setAddForm({ ...emptyForm });
                  setAddFieldErrors({});
                }}
              >
                Close
              </button>
            </div>

            <form className="modal-form" onSubmit={addCourse}>
              <div className="field">
                <label htmlFor="add-course-title">Course Title</label>
                <input
                  id="add-course-title"
                  value={addForm.title}
                  onChange={(e) => handleAddFormChange("title", e.target.value)}
                />
                {addFieldErrors.title && <p className="field-error">{addFieldErrors.title}</p>}
              </div>

              <div className="field">
                <label htmlFor="add-course-credits">Credits</label>
                <input
                  id="add-course-credits"
                  inputMode="numeric"
                  value={addForm.credits}
                  onChange={(e) => handleAddFormChange("credits", e.target.value)}
                />
                {addFieldErrors.credits && <p className="field-error">{addFieldErrors.credits}</p>}
              </div>

              <div className="field">
                <label htmlFor="add-course-code">Course Code</label>
                <input
                  id="add-course-code"
                  value={addForm.code}
                  onChange={(e) => handleAddFormChange("code", e.target.value)}
                />
                {addFieldErrors.code && <p className="field-error">{addFieldErrors.code}</p>}
              </div>

              <div className="modal-actions">
                <button
                  type="button"
                  className="btn btn-outline btn-small"
                  onClick={() => {
                    setIsAddModalOpen(false);
                    setAddForm({ ...emptyForm });
                    setAddFieldErrors({});
                  }}
                >
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary btn-small" disabled={createCourse.isPending}>
                  {createCourse.isPending ? "Saving..." : "Save Course"}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}

function EditRow({
  course,
  onDone,
  onNotify,
}: {
  course: Course;
  onDone: () => void;
  onNotify: (banner: { type: "error" | "success"; text: string }) => void;
}) {
  // Local state lets the edit row behave like a small inline form.
  const [title, setTitle] = useState(course.title);
  const [credits, setCredits] = useState(String(course.credits));
  const [code, setCode] = useState(course.code)
  const updateCourse = useUpdateCourses();

  function handleSave(e: FormEvent) {
    // Validate the edited values before sending an update request.
    e.preventDefault();
    const creditsNum = Number(credits);
    if (!title.trim() || !code.trim() || Number.isNaN(creditsNum) || creditsNum < 1 || creditsNum > 4) {
      onNotify({ type: "error", text: "Please enter a valid title, code and credits (1 - 4)." });
      return;
    }
    updateCourse.mutate(
      { courseCode: course.code, input: { title, credits: creditsNum, code } },
      {
        onSuccess: () => {
          onNotify({ type: "success", text: `Updated ${code}'s record.` });
          onDone();
        },
        onError: (err) => onNotify({ type: "error", text: apiErrorMessage(err) }),
      },
    );
  }

  return (
    <tr>
      <td colSpan={3}>
        <form className="edit-row-form" onSubmit={handleSave}>
          <input name="code" value={code} onChange={(e) => setCode(e.target.value)} aria-label="Code" />
          <input name="title" value={title} onChange={(e) => setTitle(e.target.value)} aria-label="Title" />
          <input name="credits" value={credits} onChange={(e) => setCredits(e.target.value)} aria-label="Credits" />
          
          <button type="submit" className="btn btn-primary btn-small" disabled={updateCourse.isPending}>
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
