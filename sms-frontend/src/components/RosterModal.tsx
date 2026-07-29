// RosterModal.tsx shows the list of students currently enrolled in a course,
// with their course grade, and lets you remove a student from the course.
import { toast } from "sonner";
import { Trash2 } from "lucide-react";
import Modal from "./ui/Modal";
import { Button } from "./ui/Button";
import { Badge } from "./ui/Badge";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "./ui/Table";
import { useCourseRoster, useUnenrollStudent } from "../hooks/useEnrollments";
import { apiErrorMessage } from "../lib/axios";
import type { Course } from "../types";

interface RosterModalProps {
  course: Course | null;
  onClose: () => void;
}

export function RosterModal({ course, onClose }: RosterModalProps) {
  const { data: enrollments, isLoading } = useCourseRoster(course?.code ?? null);
  const unenrollMutation = useUnenrollStudent();

  if (!course) return null;

  function handleRemove(studentId: number, name: string) {
    if (!course) return;
    if (!window.confirm(`Remove ${name} from ${course.code}?`)) return;
    unenrollMutation.mutate(
      { studentId, courseCode: course.code },
      {
        onSuccess: () => toast.success(`Removed ${name} from ${course.code}.`),
        onError: (err) => toast.error(apiErrorMessage(err)),
      },
    );
  }

  return (
    <Modal isOpen={!!course} onClose={onClose} title={`${course.title} Roster`}>
      <div className="space-y-4">
        {isLoading && <p className="text-sm text-muted-foreground">Loading roster…</p>}

        {!isLoading && !enrollments?.length && (
          <p className="text-sm text-muted-foreground">No students enrolled in this course yet.</p>
        )}

        {!isLoading && !!enrollments?.length && (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Student ID</TableHead>
                <TableHead>Name</TableHead>
                <TableHead>Grade Level</TableHead>
                <TableHead>Course Grade</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {enrollments.map((enrollment) => (
                <TableRow key={enrollment.studentId}>
                  <TableCell>
                    <Badge variant="outline" className="font-mono">
                      #{enrollment.studentId}
                    </Badge>
                  </TableCell>
                  <TableCell className="font-medium">{enrollment.student?.name ?? "—"}</TableCell>
                  <TableCell>{enrollment.student?.grade ?? "—"}</TableCell>
                  <TableCell>{enrollment.grade || "—"}</TableCell>
                  <TableCell>
                    <div className="flex justify-end">
                      <Button
                        variant="destructive"
                        size="sm"
                        onClick={() =>
                          handleRemove(enrollment.studentId, enrollment.student?.name ?? "this student")
                        }
                        disabled={unenrollMutation.isPending}
                      >
                        <Trash2 />
                        Remove
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </Modal>
  );
}
