// StudentsView.tsx renders the roster UI and handles filtering, editing, and deletion.
import { useState, type FormEvent } from "react";
import { Plus, Pencil, Trash2, Check, X } from "lucide-react";
import { toast } from "sonner";
import { useDeleteStudent, useStudents, useUpdateStudent, useCreateStudent, type StudentFilter } from "../hooks/useStudents";
import { apiErrorMessage } from "../lib/axios";
import { StudentInputSchema, type Student, type StudentInput } from "../types";
import PageHeader from "./ui/PageHeader";
import { Card } from "./ui/Card";
import { Button } from "./ui/Button";
import Input from "./ui/Input";
import Modal from "./ui/Modal";
import { Badge } from "./ui/Badge";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "./ui/Select";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "./ui/Table";

type FilterMode = "all" | "grade" | "honors" | "search" | "name";

export function StudentsView() {
  const emptyForm = { studentId: "", name: "", grade: "", gpa: "" };

  const [mode, setMode] = useState<FilterMode>("all");
  const [gradeInput, setGradeInput] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [nameInput, setNameInput] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [addForm, setAddForm] = useState(emptyForm);
  const [addFieldErrors, setAddFieldErrors] = useState<Partial<Record<keyof StudentInput, string>>>({});

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

  const { data: students, isLoading, isError, error } = useStudents(filter);
  const deleteStudent = useDeleteStudent();
  const createStudent = useCreateStudent();

  function handleDelete(student: Student) {
    if (!window.confirm(`Remove ${student.name} (ID ${student.studentId}) from the roster`)) return;
    deleteStudent.mutate(student.studentId, {
      onSuccess: () => toast.success(`Removed ${student.name} from the roster.`),
      onError: (err) => toast.error(apiErrorMessage(err)),
    });
  }

  function handleAddFormChange(field: keyof typeof addForm, value: string) {
    setAddForm((prev) => ({ ...prev, [field]: value }));
  }

  function openAddStudentModal() {
    setAddForm({ ...emptyForm });
    setAddFieldErrors({});
    setIsAddModalOpen(true);
  }

  function closeAddStudentModal() {
    setIsAddModalOpen(false);
    setAddForm({ ...emptyForm });
    setAddFieldErrors({});
  }

  function addStudent(e: FormEvent) {
    e.preventDefault();

    const result = StudentInputSchema.safeParse(addForm);
    if (!result.success) {
      const errors: Partial<Record<keyof StudentInput, string>> = {};
      for (const issue of result.error.issues) {
        const key = issue.path[0] as keyof StudentInput;
        errors[key] = issue.message;
      }
      setAddFieldErrors(errors);
      return;
    }

    setAddFieldErrors({});
    createStudent.mutate(result.data, {
      onSuccess: () => {
        setAddForm({ ...emptyForm });
        setIsAddModalOpen(false);
        toast.success(`Enrolled ${result.data.name} (ID ${result.data.studentId}) in the roster.`);
      },
      onError: (err) => toast.error(apiErrorMessage(err)),
    });
  }

  return (
    <div>
      <PageHeader
        title="Student Roster"
        description="View, filter, edit, add or remove students on file."
        actions={
          <Button onClick={openAddStudentModal}>
            <Plus />
            Enroll a Student
          </Button>
        }
      />

      <Card padding="responsive" className="mb-6">
        <div className="flex flex-wrap items-end gap-4">
          <div className="w-full sm:w-56">
            <label htmlFor="filter-mode" className="block text-sm font-medium text-foreground/80 mb-2">
              Show
            </label>
            <Select
              value={mode}
              onValueChange={(value) => setMode(value as FilterMode)}
            >
              <SelectTrigger id="filter-mode" className="w-full">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All students</SelectItem>
                <SelectItem value="grade">By grade</SelectItem>
                <SelectItem value="honors">Honors (GPA ≥ 3.5)</SelectItem>
                <SelectItem value="search">Search by student ID</SelectItem>
                <SelectItem value="name">Search by name</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {mode === "grade" && (
            <div className="w-full sm:w-40">
              <Input
                label="Grade"
                id="grade-filter"
                inputMode="numeric"
                placeholder="3"
                value={gradeInput}
                onChange={(e) => setGradeInput(e.target.value)}
              />
            </div>
          )}

          {mode === "search" && (
            <div className="w-full sm:w-48">
              <Input
                label="Student ID"
                id="search-filter"
                inputMode="numeric"
                placeholder="1001"
                value={searchInput}
                onChange={(e) => setSearchInput(e.target.value)}
              />
            </div>
          )}

          {mode === "name" && (
            <div className="w-full sm:w-64">
              <Input
                label="Student Name"
                id="name-filter"
                inputMode="text"
                placeholder="Khaymar Moe"
                value={nameInput}
                onChange={(e) => setNameInput(e.target.value)}
              />
            </div>
          )}
        </div>
      </Card>

      {isLoading && <p className="text-sm text-muted-foreground">Loading roster…</p>}
      {isError && (
        <div className="mb-4 rounded-lg border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {apiErrorMessage(error)}
        </div>
      )}

      {!isLoading && !isError && !students?.length && (
        <Card padding="lg" className="text-center text-muted-foreground">
          No students match this view yet.
        </Card>
      )}

      {!isLoading && !!students?.length && (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ID</TableHead>
              <TableHead>Name</TableHead>
              <TableHead>Grade</TableHead>
              <TableHead>GPA</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {students.map((student) =>
              editingId === student.studentId ? (
                <EditRow key={student.studentId} student={student} onDone={() => setEditingId(null)} />
              ) : (
                <TableRow key={student.studentId}>
                  <TableCell>
                    <Badge variant="outline" className="font-mono">
                      #{student.studentId}
                    </Badge>
                  </TableCell>
                  <TableCell className="font-medium">{student.name}</TableCell>
                  <TableCell>{student.grade}</TableCell>
                  <TableCell>
                    <span className={student.gpa >= 3.5 ? "text-success font-semibold" : undefined}>
                      {student.gpa.toFixed(2)}
                    </span>
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center justify-end gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setEditingId(student.studentId)}
                      >
                        <Pencil />
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDelete(student)}
                        disabled={deleteStudent.isPending}
                      >
                        <Trash2 />
                        Remove
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ),
            )}
          </TableBody>
        </Table>
      )}

      <Modal isOpen={isAddModalOpen} onClose={closeAddStudentModal} title="Enroll a Student">
        <form className="space-y-4" onSubmit={addStudent}>
          <Input
            label="Student ID"
            id="studentId"
            inputMode="numeric"
            placeholder="1005"
            value={addForm.studentId}
            onChange={(e) => handleAddFormChange("studentId", e.target.value)}
            error={addFieldErrors.studentId}
            required
          />

          <Input
            label="Full Name"
            id="name"
            placeholder="Khaymar Moe"
            value={addForm.name}
            onChange={(e) => handleAddFormChange("name", e.target.value)}
            error={addFieldErrors.name}
            required
          />

          <Input
            label="Grade"
            id="grade"
            inputMode="numeric"
            placeholder="3"
            value={addForm.grade}
            onChange={(e) => handleAddFormChange("grade", e.target.value)}
            error={addFieldErrors.grade}
            required
          />

          <Input
            label="GPA"
            id="gpa"
            inputMode="decimal"
            placeholder="3.8"
            value={addForm.gpa}
            onChange={(e) => handleAddFormChange("gpa", e.target.value)}
            error={addFieldErrors.gpa}
            required
          />

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={closeAddStudentModal}>
              Cancel
            </Button>
            <Button type="submit" disabled={createStudent.isPending}>
              {createStudent.isPending ? "Saving..." : "Save Student"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

function EditRow({
  student,
  onDone,
}: {
  student: Student;
  onDone: () => void;
}) {
  const [name, setName] = useState(student.name);
  const [grade, setGrade] = useState(String(student.grade));
  const [gpa, setGpa] = useState(String(student.gpa));
  const updateStudent = useUpdateStudent();

  function handleSave(e: FormEvent) {
    e.preventDefault();
    const gpaNum = Number(gpa);
    const gradeNum = Number(grade);
    if (!name.trim() || Number.isNaN(gpaNum) || Number.isNaN(gradeNum) || gpaNum < 0 || gpaNum > 4) {
      toast.error("Please enter a valid name, grade, and GPA (0.0–4.0).");
      return;
    }
    updateStudent.mutate(
      { studentId: student.studentId, input: { studentId: student.studentId, name, grade: gradeNum, gpa: gpaNum } },
      {
        onSuccess: () => {
          toast.success(`Updated ${name}'s record.`);
          onDone();
        },
        onError: (err) => toast.error(apiErrorMessage(err)),
      },
    );
  }

  return (
    <TableRow>
      <TableCell>
        <Badge variant="outline" className="font-mono">
          #{student.studentId}
        </Badge>
      </TableCell>
      <TableCell colSpan={3}>
        <form className="flex flex-wrap items-center gap-2" onSubmit={handleSave}>
          <Input
            className="max-w-[160px]"
            name="name"
            inputSize="sm"
            value={name}
            onChange={(e) => setName(e.target.value)}
            aria-label="Name"
          />
          <Input
            className="max-w-[90px]"
            name="grade"
            inputSize="sm"
            value={grade}
            onChange={(e) => setGrade(e.target.value)}
            aria-label="Grade"
          />
          <Input
            className="max-w-[90px]"
            name="gpa"
            inputSize="sm"
            value={gpa}
            onChange={(e) => setGpa(e.target.value)}
            aria-label="GPA"
          />
          <Button type="submit" size="sm" disabled={updateStudent.isPending}>
            <Check />
            Save
          </Button>
          <Button type="button" variant="outline" size="sm" onClick={onDone}>
            <X />
            Cancel
          </Button>
        </form>
      </TableCell>
      <TableCell />
    </TableRow>
  );
}
