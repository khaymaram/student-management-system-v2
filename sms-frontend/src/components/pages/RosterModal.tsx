// RosterModal.tsx shows the list of students currently enrolled in a course,
// with their course grade, and lets you remove a student from the course.
import { toast } from "sonner";
import { Trash2, GraduationCap} from "lucide-react";
import Modal from "../ui/Modal";
import { Button } from "../ui/Button";
import { Badge } from "../ui/Badge";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "../ui/Table";
import { useCourseRoster, useUnenrollStudent, useEnrollStudent } from "../../hooks/useEnrollments"
import { apiErrorMessage } from "../../lib/axios";
import type { Course } from "../../types";
import { useState, useMemo } from "react";
import { useStudents } from "../../hooks/useStudents";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "../ui/Select";

interface RosterModalProps {
    course: Course | null;
    onClose: () => void;
}

export function RosterModal({ course, onClose }: RosterModalProps) {
    const [selectedStudent, setSelectedStudent] = useState<string>();

    const { data: enrollments, isLoading } = useCourseRoster(course?.code ?? null);
    const { data: allStudents } = useStudents({ type: "all" });
    const enrollMutation = useEnrollStudent();
    const unenrollMutation = useUnenrollStudent();

    const enrolledIds = useMemo(
        () => new Set((enrollments ?? []).map((e) => e.studentId)),
        [enrollments],
    )

    const availableStudents = (allStudents ?? []).filter(
        (student) => !enrolledIds.has(student.studentId)
    );

    if (!course) return null;

    function handleEnroll() {
        if (!course || !selectedStudent) return;
        // if (isOverCreditLimit) {
        //     toast.error("This enrollment would push the student over 15 credits.");
        //     return;
        // }

        enrollMutation.mutate(
            { courseCode: course.code, studentId: Number(selectedStudent) },
            {
                onSuccess: () => {
                    toast.success(`Enrolled ${selectedStudent} in ${course.title}.`);
                    setSelectedStudent(undefined);
                },
                onError: (err) => toast.error(apiErrorMessage(err)),
            },
        );
    }

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

                <div className="flex flex-wrap items-end gap-3">
                    <div className="flex-1 min-w-[180px]">
                        <label className="block text-sm font-medium text-foreground/80 mb-2">
                            Enroll a Student
                        </label>
                        <Select value={selectedStudent} onValueChange={setSelectedStudent}>
                            <SelectTrigger className="w-full">
                                <SelectValue placeholder="Select a student" />
                            </SelectTrigger>
                            <SelectContent>
                                {availableStudents.length === 0 ? (
                                    <div className="px-3 py-2 text-sm text-muted-foreground">
                                        {"No available students"}
                                    </div>
                                ) : (
                                    availableStudents.map((student) => (
                                        <SelectItem key={student.studentId} value={String(student.studentId)}>
                                            {student.name}
                                        </SelectItem>
                                    ))
                                )}
                            </SelectContent>
                        </Select>
                    </div>
                    <Button
                        onClick={handleEnroll}
                        disabled={!selectedStudent || enrollMutation.isPending}
                    >
                        <GraduationCap />
                        Enroll
                    </Button>
                </div>

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