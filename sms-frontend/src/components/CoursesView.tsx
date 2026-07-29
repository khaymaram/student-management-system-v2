import { useState, type FormEvent } from "react";
import { Plus, Pencil, Trash2, Check, X } from "lucide-react";
import { toast } from "sonner";
import { useDeleteCourse, useCourses, useUpdateCourses, useCreateCourse, type CourseFilter } from "../hooks/useCourses";
import { apiErrorMessage } from "../lib/axios";
import { CourseInputSchema, type Course, type CourseInput } from "../types";
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

type FilterMode = "all" | "credits" | "title" | "code";

export function CoursesView() {
  const emptyForm = { title: "", credits: "", code: "" };
  const [mode, setMode] = useState<FilterMode>("all");
  const [creditInput, setCreditInput] = useState("");
  const [titleInput, setTitleInput] = useState("");
  const [codeInput, setCodeInput] = useState("");
  const [editingCode, setEditingCode] = useState<string | null>(null);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [addForm, setAddForm] = useState(emptyForm);
  const [addFieldErrors, setAddFieldErrors] = useState<Partial<Record<keyof CourseInput, string>>>({});

  const filter: CourseFilter =
    mode === "credits" && creditInput
      ? { type: "credits", credits: Number(creditInput) }
      : mode === "title" && titleInput
        ? { type: "title", title: String(titleInput) }
        : mode === "code" && codeInput
          ? { type: "code", code: String(codeInput) }
          : { type: "all" };

  const { data: courses, isLoading, isError, error } = useCourses(filter);
  const deleteCourse = useDeleteCourse();
  const createCourse = useCreateCourse();

  function handleDelete(course: Course) {
    if (!window.confirm(`Remove ${course.title} (${course.code}) from the directory`)) return;
    deleteCourse.mutate(course.code, {
      onSuccess: () => toast.success(`Removed ${course.title} from the directory.`),
      onError: (err) => toast.error(apiErrorMessage(err)),
    });
  }

  function handleAddFormChange(field: keyof typeof addForm, value: string) {
    setAddForm((prev) => ({ ...prev, [field]: value }));
  }

  function openAddCourseModal() {
    setAddForm({ ...emptyForm });
    setAddFieldErrors({});
    setIsAddModalOpen(true);
  }

  function closeAddCourseModal() {
    setIsAddModalOpen(false);
    setAddForm({ ...emptyForm });
    setAddFieldErrors({});
  }

  function addCourse(e: FormEvent) {
    e.preventDefault();

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
        toast.success(`Created ${result.data.title} (${result.data.code}) in the directory.`);
      },
      onError: (err) => toast.error(apiErrorMessage(err)),
    });
  }

  return (
    <div>
      <PageHeader
        title="Course Directory"
        description="View, filter, edit, add or remove courses."
        actions={
          <Button onClick={openAddCourseModal}>
            <Plus />
            Add a Course
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
                <SelectItem value="all">All courses</SelectItem>
                <SelectItem value="credits">By credits</SelectItem>
                <SelectItem value="title">Search by course title</SelectItem>
                <SelectItem value="code">Search by course code</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {mode === "credits" && (
            <div className="w-full sm:w-40">
              <Input
                label="Credits"
                id="credits-filter"
                inputMode="numeric"
                placeholder="3"
                value={creditInput}
                onChange={(e) => setCreditInput(e.target.value)}
              />
            </div>
          )}

          {mode === "title" && (
            <div className="w-full sm:w-64">
              <Input
                label="Course Title"
                id="title-filter"
                inputMode="text"
                placeholder="Intro to OOP"
                value={titleInput}
                onChange={(e) => setTitleInput(e.target.value)}
              />
            </div>
          )}

          {mode === "code" && (
            <div className="w-full sm:w-48">
              <Input
                label="Course Code"
                id="code-filter"
                inputMode="text"
                placeholder="CMSC131"
                value={codeInput}
                onChange={(e) => setCodeInput(e.target.value)}
              />
            </div>
          )}
        </div>
      </Card>

      {isLoading && <p className="text-sm text-muted-foreground">Loading courses…</p>}
      {isError && (
        <div className="mb-4 rounded-lg border border-destructive/40 bg-destructive/10 px-4 py-3 text-sm text-destructive">
          {apiErrorMessage(error)}
        </div>
      )}

      {!isLoading && !isError && !courses?.length && (
        <Card padding="lg" className="text-center text-muted-foreground">
          No courses match this view yet.
        </Card>
      )}

      {!isLoading && !!courses?.length && (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Course Code</TableHead>
              <TableHead>Title</TableHead>
              <TableHead>Credits</TableHead>
              <TableHead className="text-right">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {courses.map((course) =>
              editingCode === course.code ? (
                <EditRow key={course.code} course={course} onDone={() => setEditingCode(null)} />
              ) : (
                <TableRow key={course.code}>
                  <TableCell>
                    <Badge variant="outline" className="font-mono">
                      {course.code}
                    </Badge>
                  </TableCell>
                  <TableCell className="font-medium">{course.title}</TableCell>
                  <TableCell>{course.credits}</TableCell>
                  <TableCell>
                    <div className="flex items-center justify-end gap-2">
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setEditingCode(course.code)}
                      >
                        <Pencil />
                        Edit
                      </Button>
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() => handleDelete(course)}
                        disabled={deleteCourse.isPending}
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

      <Modal isOpen={isAddModalOpen} onClose={closeAddCourseModal} title="Add a Course">
        <form className="space-y-4" onSubmit={addCourse}>
          <Input
            label="Course Title"
            id="add-course-title"
            value={addForm.title}
            onChange={(e) => handleAddFormChange("title", e.target.value)}
            error={addFieldErrors.title}
            required
          />

          <Input
            label="Credits"
            id="add-course-credits"
            inputMode="numeric"
            value={addForm.credits}
            onChange={(e) => handleAddFormChange("credits", e.target.value)}
            error={addFieldErrors.credits}
            required
          />

          <Input
            label="Course Code"
            id="add-course-code"
            value={addForm.code}
            onChange={(e) => handleAddFormChange("code", e.target.value)}
            error={addFieldErrors.code}
            required
          />

          <div className="flex justify-end gap-2 pt-2">
            <Button type="button" variant="outline" onClick={closeAddCourseModal}>
              Cancel
            </Button>
            <Button type="submit" disabled={createCourse.isPending}>
              {createCourse.isPending ? "Saving..." : "Save Course"}
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}

function EditRow({
  course,
  onDone,
}: {
  course: Course;
  onDone: () => void;
}) {
  const [title, setTitle] = useState(course.title);
  const [credits, setCredits] = useState(String(course.credits));
  const [code, setCode] = useState(course.code);
  const updateCourse = useUpdateCourses();

  function handleSave(e: FormEvent) {
    e.preventDefault();
    const creditsNum = Number(credits);
    if (!title.trim() || !code.trim() || Number.isNaN(creditsNum) || creditsNum < 1 || creditsNum > 4) {
      toast.error("Please enter a valid title, code and credits (1 - 4).");
      return;
    }
    updateCourse.mutate(
      { courseCode: course.code, input: { title, credits: creditsNum, code } },
      {
        onSuccess: () => {
          toast.success(`Updated ${code}'s record.`);
          onDone();
        },
        onError: (err) => toast.error(apiErrorMessage(err)),
      },
    );
  }

  return (
    <TableRow>
      <TableCell colSpan={3}>
        <form className="flex flex-wrap items-center gap-2" onSubmit={handleSave}>
          <Input
            className="max-w-[120px]"
            name="code"
            inputSize="sm"
            value={code}
            onChange={(e) => setCode(e.target.value)}
            aria-label="Code"
          />
          <Input
            className="max-w-[200px]"
            name="title"
            inputSize="sm"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            aria-label="Title"
          />
          <Input
            className="max-w-[90px]"
            name="credits"
            inputSize="sm"
            value={credits}
            onChange={(e) => setCredits(e.target.value)}
            aria-label="Credits"
          />
          <Button type="submit" size="sm" disabled={updateCourse.isPending}>
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
