
// StudentsView.tsx renders the student roster with server-side pagination,
// filtering, editing, creation, and deletion.

import {
  useEffect,
  useState,
  type FormEvent,
} from "react";

import {
  BookOpen,
  Pencil,
  Plus,
  Trash2,
} from "lucide-react";

import { Link } from "react-router-dom";
import { toast } from "sonner";

import {
  useCreateStudent,
  useDeleteStudent,
  useStudentsPaginated,
  useUpdateStudent,
  type StudentFilter,
} from "../../hooks/useStudents";

import { apiErrorMessage } from "../../lib/axios";

import {
  StudentInputSchema,
  type Student,
  type StudentInput,
} from "../../types";

import PageHeader from "../ui/PageHeader";
import { Card } from "../ui/Card";
import { Button } from "../ui/Button";
import Input from "../ui/Input";
import Modal from "../ui/Modal";
import { Badge } from "../ui/Badge";
import { Pagination } from "../ui/Pagination";

import { EnrollmentModal } from "./EnrollmentModal";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "../ui/Select";

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/Table";

type FilterMode =
  | "all"
  | "grade"
  | "honors"
  | "search"
  | "name";

const emptyForm = {
  studentId: "",
  name: "",
  grade: "",
  gpa: "",
};

export function StudentsView() {

  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(5);

  const [mode, setMode] = useState<FilterMode>("all");
  const [gradeInput, setGradeInput] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [nameInput, setNameInput] = useState("");
  const [enrollmentStudent, setEnrollmentStudent,] = useState<Student | null>(null);
  const [editingStudent, setEditingStudent,] = useState<Student | null>(null);

  const [isAddModalOpen, setIsAddModalOpen,] = useState(false);
  const [addForm, setAddForm] = useState(emptyForm);

  const [addFieldErrors, setAddFieldErrors,] = useState<Partial<Record<keyof StudentInput, string>>>({});

  const [editForm, setEditForm] = useState(emptyForm);

  const [editFieldErrors, setEditFieldErrors,] = useState<Partial<Record<keyof StudentInput, string>>>({});

  const filter: StudentFilter =
    mode === "grade" && gradeInput ? {
      type: "grade",
      grade: Number(gradeInput),
    } : mode === "honors" ? {
      type: "honors",
    } : mode === "search" && searchInput ? {
      type: "search",
      studentId: Number(searchInput),
    } : mode === "name" && nameInput.trim() ? {
      type: "name",
      name: nameInput.trim(),
    } : {
      type: "all",
    };

  // Return to the first page whenever
  // the active filter changes.
  useEffect(() => {
    setPage(1);
  }, [mode, gradeInput, searchInput, nameInput,]);

  const { data, isLoading, isFetching, isError, error,
  } = useStudentsPaginated({ page, pageSize, filter, });

  const students = data?.data ?? [];

  const totalPages = data?.totalPages ?? 1;

  const totalCount = data?.total ?? 0;

  const deleteStudent = useDeleteStudent();

  const createStudent = useCreateStudent();

  const updateStudent = useUpdateStudent();

  function handleDelete(student: Student) {
    if (!window.confirm(`Remove ${student.name} (ID ${student.studentId}) from the roster?`)) {
      return;
    }
    deleteStudent.mutate(student.studentId, {
      onSuccess: () => toast.success(`Removed ${student.name} from the roster.`),
      onError: (err) => toast.error(apiErrorMessage(err)),
    });
  }

  function handleAddFormChange(field: keyof typeof addForm, value: string) {
    setAddForm((prev) => ({ ...prev, [field]: value }));
  }

  function handleEditFormChange(
    field: keyof typeof editForm,
    value: string
  ) {
    setEditForm((prev) => ({
      ...prev,
      [field]: value,
    }));
  }

  function openEditStudentModal(student: Student) {
    setEditingStudent(student);

    setEditForm({
      studentId: String(student.studentId), name: student.name, grade: String(student.grade), gpa: String(student.gpa),
    });

    setEditFieldErrors({});
  }

  function closeEditStudentModal() {
    setEditingStudent(null);
    setEditForm({
      ...emptyForm,
    });
    setEditFieldErrors({});
  }

  function openAddStudentModal() {
    setAddForm({
      ...emptyForm,
    });

    setAddFieldErrors({});
    setIsAddModalOpen(true);
  }

  function closeAddStudentModal() {
    setIsAddModalOpen(false);

    setAddForm({
      ...emptyForm,
    });

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

  function editStudent(e: FormEvent) {
    e.preventDefault();

    if (!editingStudent) return;

    const result = StudentInputSchema.safeParse(editForm);

    if (!result.success) {

      const errors: Partial<Record<keyof StudentInput, string>> = {};

      for (const issue of result.error.issues) {
        const key = issue.path[0] as keyof StudentInput;
        errors[key] = issue.message;
      }

      setEditFieldErrors(errors);

      return;
    }

    setEditFieldErrors({});

    updateStudent.mutate(
      {
        studentId: editingStudent.studentId,
        input: result.data,
      },
      {
        onSuccess: () => {
          toast.success(`Updated ${result.data.name}.`);
          closeEditStudentModal();
        },
        onError: (err) => toast.error(apiErrorMessage(err)),
      }
    );
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

      {!isLoading &&!isError &&students.length === 0 && (
          <Card padding="lg" className="text-center text-muted-foreground">
            No students found.
          </Card>
        )}

      {!isLoading &&!isError &&students.length > 0 && (
          <>
            <div className={isFetching  ? "opacity-70 transition-opacity"  : undefined}>

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

                  {students.map((student) => (
                      <TableRow
                        key={student.studentId}
                      >

                        <TableCell>

                          <Link to={`/roster/${student.studentId}`}>
                            <Badge variant="outline"
                              className="cursor-pointer font-mono transition-colors hover:bg-accent hover:text-accent-foreground"
                            >
                              #{student.studentId}
                            </Badge>
                          </Link>

                        </TableCell>

                        <TableCell className="font-medium">{student.name}</TableCell>

                        <TableCell>{student.grade}</TableCell>

                        <TableCell>
                          <span className={student.gpa >= 3.5? "font-bold text-success": undefined}>
                            {student.gpa.toFixed(2)}
                          </span>

                        </TableCell>

                        <TableCell>
                          <div className="flex items-center justify-end gap-2">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() =>setEnrollmentStudent(student)}
                            >
                              <BookOpen />
                              Courses
                            </Button>

                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() =>openEditStudentModal(student)}
                            >
                              <Pencil />
                              Edit
                            </Button>

                            <Button
                              variant="destructive"
                              size="sm"
                              onClick={() =>handleDelete(student)
                              }
                              disabled={deleteStudent.isPending}
                            >
                              <Trash2 />
                              Remove
                            </Button>
                          </div>
                        </TableCell>
                      </TableRow>
                    )
                  )}
                </TableBody>
              </Table>
            </div>

            <Pagination
              page={page}
              totalPages={totalPages}
              pageSize={pageSize}
              totalCount={totalCount}
              onPageChange={setPage}
              onPageSizeChange={(size) => {setPageSize(size);setPage(1);}}
            />

          </>
        )}

      <Modal isOpen={isAddModalOpen} onClose={closeAddStudentModal} title="Enroll a Student">
        <form className="space-y-4" onSubmit={addStudent}>

          <Input
            label="Student ID"
            id="studentId"
            inputMode="numeric"
            placeholder="1005"
            value={
              addForm.studentId
            }
            onChange={(e) =>
              handleAddFormChange(
                "studentId",
                e.target.value
              )
            }
            error={
              addFieldErrors.studentId
            }
            required
          />

          <Input
            label="Full Name"
            id="name"
            placeholder="Khaymar Moe"
            value={
              addForm.name
            }
            onChange={(e) =>
              handleAddFormChange(
                "name",
                e.target.value
              )
            }
            error={
              addFieldErrors.name
            }
            required
          />

          <Input
            label="Grade"
            id="grade"
            inputMode="numeric"
            placeholder="3"
            value={
              addForm.grade
            }
            onChange={(e) =>
              handleAddFormChange(
                "grade",
                e.target.value
              )
            }
            error={
              addFieldErrors.grade
            }
            required
          />

          <Input
            label="GPA"
            id="gpa"
            inputMode="decimal"
            placeholder="3.8"
            value={
              addForm.gpa
            }
            onChange={(e) =>
              handleAddFormChange(
                "gpa",
                e.target.value
              )
            }
            error={
              addFieldErrors.gpa
            }
            required
          />

          <div className="flex justify-end gap-2 pt-2">

            <Button
              type="button"
              variant="outline"
              onClick={
                closeAddStudentModal
              }
            >
              Cancel
            </Button>

            <Button
              type="submit"
              disabled={
                createStudent.isPending
              }
            >
              {createStudent.isPending
                ? "Saving..."
                : "Save Student"}
            </Button>

          </div>

        </form>
      </Modal>

      <Modal
        isOpen={
          !!editingStudent
        }
        onClose={
          closeEditStudentModal
        }
        title="Edit Student"
      >
        <form
          className="space-y-4"
          onSubmit={
            editStudent
          }
        >

          <Input
            label="Student ID"
            value={
              editForm.studentId
            }
            disabled
          />

          <Input
            label="Full Name"
            value={
              editForm.name
            }
            onChange={(e) =>
              handleEditFormChange(
                "name",
                e.target.value
              )
            }
            error={
              editFieldErrors.name
            }
            required
          />

          <Input
            label="Grade"
            inputMode="numeric"
            value={
              editForm.grade
            }
            onChange={(e) =>
              handleEditFormChange(
                "grade",
                e.target.value
              )
            }
            error={
              editFieldErrors.grade
            }
            required
          />

          <Input
            label="GPA"
            inputMode="decimal"
            value={
              editForm.gpa
            }
            onChange={(e) =>
              handleEditFormChange(
                "gpa",
                e.target.value
              )
            }
            error={
              editFieldErrors.gpa
            }
            required
          />

          <div className="flex justify-end gap-2 pt-2">

            <Button
              type="button"
              variant="outline"
              onClick={
                closeEditStudentModal
              }
            >
              Cancel
            </Button>

            <Button
              type="submit"
              disabled={
                updateStudent.isPending
              }
            >
              {updateStudent.isPending
                ? "Saving..."
                : "Save Changes"}
            </Button>

          </div>

        </form>
      </Modal>

      <EnrollmentModal student={enrollmentStudent}onClose={() =>setEnrollmentStudent(null)}/>
    </div>
  );
}