// page details will have: 
//  back to students button at the top of the page 
// Name
// in the header area there will be: number of courses enrolled, grade level, credits, gpa
// body: table with courses enrolled: id, title, course grade, unenroll from course button


import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { ArrowLeft, Check, Trash2, Plus } from "lucide-react";
import { toast } from "sonner";
import { EnrollmentModal } from "./EnrollmentModal";

import { Button } from "../ui/Button";
import { Badge } from "../ui/Badge";
import {
    Table,
    TableHeader,
    TableBody,
    TableRow,
    TableHead,
    TableCell,
} from "../ui/Table";
import { Card } from "../ui/Card";

import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "../ui/Select";

import {
    useStudentEnrollments,
    useUpdateEnrollmentGrade,
    useUnenrollStudent,
} from "../../hooks/useEnrollments";

import { useStudent } from "../../hooks/useStudents";

import { apiErrorMessage } from "../../lib/axios";

import type { Student, Enrollment } from "../../types";

const GRADE_OPTIONS = ["A", "B", "C", "D", "F"] as const;

function getTotalCredits(enrollments: Enrollment[] | undefined) {
    return (
        enrollments?.reduce(
            (sum, enrollment) => sum + (enrollment.course?.credits ?? 0),
            0
        ) ?? 0
    );
}

export default function StudentDetailsView() {
    const { studentId } = useParams();
    const navigate = useNavigate();
    const parsedStudentId = studentId ? Number(studentId) : null;
    const [enrollmentStudent, setEnrollmentStudent] = useState<Student | null>(null);

    const {
        data: studentData,
        isLoading: studentLoading,
    } = useStudent(parsedStudentId);

    const student = studentData as Student | undefined;

    const {
        data: enrollments,
        isLoading: enrollmentsLoading,
    } = useStudentEnrollments(parsedStudentId);

    const totalCredits = useMemo(
        () => getTotalCredits(enrollments),
        [enrollments]
    );

    // const averageGrade = useMemo(() => {
    //     if (!enrollments?.length) return null;

    //     const points: Record<string, number> = {
    //         A: 4,
    //         B: 3,
    //         C: 2,
    //         D: 1,
    //         F: 0,
    //     };

    //     const graded = enrollments.filter(
    //         (e) => e.grade && points[e.grade] !== undefined
    //     );

    //     if (graded.length === 0) return null;

    //     const total = graded.reduce(
    //         (sum, e) => sum + points[e.grade!],
    //         0
    //     );

    //     const avg = total / graded.length;

    //     if (avg >= 3.5) return "A";
    //     if (avg >= 2.5) return "B";
    //     if (avg >= 1.5) return "C";
    //     if (avg >= 0.5) return "D";
    //     return "F";


    //     // return (total / graded.length).toFixed(2);
    // }, [enrollments]);

    if (!studentId) {
        return (
            <div className="p-6">
                <p>Student not found.</p>
            </div>
        );
    }

    return (
        <div className="space-y-6">

            {/* Back Button */}
            <Button
                variant="outline"
                onClick={() => navigate("/roster")}
                className="w-fit"
            >
                <ArrowLeft />
                Back to Roster
            </Button>

            {/* Course Header */}
            {!studentLoading && student && (
                <div className="space-y-1">

                    <h1 className="text-4xl font-bold tracking-tight">
                        {student.name}
                    </h1>

                    <p className="text-lg text-muted-foreground font-mono">
                        Student ID: {student.studentId}
                    </p>

                </div>
            )}

            {/* Statistics */}
            {!studentLoading && student && (
                <Card padding="responsive">

                    <div className="grid grid-cols-1 gap-6 md:grid-cols-4">

                        <div>
                            <p className="text-sm text-muted-foreground">
                                Courses Enrolled
                            </p>

                            <p className="mt-1 text-3xl font-bold">
                                {enrollments?.length ?? 0}
                            </p>
                        </div>
                        <div>
                            <p className="text-sm text-muted-foreground">
                                Grade
                            </p>

                            <p className="mt-1 text-3xl font-bold">

                                {student.grade}

                            </p>
                        </div>
                        <div>
                            <p className="text-sm text-muted-foreground">
                                Credits
                            </p>

                            <p className="mt-1 text-3xl font-bold">
                                {totalCredits}
                            </p>
                        </div>


                        <div>
                            <p className="text-sm text-muted-foreground">
                                GPA
                            </p>

                            <p className="mt-1 text-3xl font-bold">
                                {student.gpa.toFixed(2)}
                            </p>
                        </div>

                    </div>

                </Card>
            )}

            {/* Loading */}
            {enrollmentsLoading && (
                <Card padding="lg">
                    <p className="text-muted-foreground">
                        Loading enrollments...
                    </p>
                </Card>
            )}

            {/* Empty */}
            {!enrollmentsLoading && student &&
                (!enrollments || enrollments.length === 0) && (
                    <Card padding="lg">

                        <div className="mb-6 grid grid-cols-1 gap-6 md:grid-cols-3">
                            <div></div>
                            <p className="text-muted-foreground center-align">
                                {student.name} is not enrolled in any courses.
                            </p>

                            <div></div>
                            <div></div>
                            <Button onClick={() => setEnrollmentStudent(student)}>
                                <Plus />
                                Add Course
                            </Button>
                        </div>

                    </Card>
                )}

            {/* Roster */}
            {!enrollmentsLoading &&
                enrollments && student &&
                enrollments.length > 0 && (

                    <Card padding="responsive">

                        <div className="mb-6 grid grid-cols-1 gap-6 md:grid-cols-4">

                            <h2 className="text-xl font-semibold">
                                {student.name}'s Courses
                            </h2>

                            <div></div>
                            <div></div>

                            <Button onClick={() => setEnrollmentStudent(student)}>
                                <Plus />
                                Add Course
                            </Button>

                        </div>

                        <Table>

                            <TableHeader>

                                <TableRow>

                                    <TableHead>
                                        Course ID
                                    </TableHead>

                                    <TableHead>
                                        Title
                                    </TableHead>

                                    <TableHead>
                                        Professor
                                    </TableHead>

                                    <TableHead>
                                        Credits
                                    </TableHead>

                                    <TableHead>
                                        Course Grade
                                    </TableHead>

                                    <TableHead className="text-right">
                                        Actions
                                    </TableHead>

                                </TableRow>

                            </TableHeader>

                            <TableBody>

                                {enrollments.map((enrollment) => (
                                    <RosterRow
                                        key={enrollment.courseCode}
                                        enrollment={enrollment}
                                    />
                                ))}

                            </TableBody>

                        </Table>

                    </Card>

                )}
            <EnrollmentModal student={enrollmentStudent} onClose={() => setEnrollmentStudent(null)} />
        </div>
    );
}

function RosterRow({
    enrollment,
}: {
    enrollment: Enrollment;
}) {
    const [gradeInput, setGradeInput] = useState(
        enrollment.grade ?? ""
    );

    const updateGrade = useUpdateEnrollmentGrade();
    const unenrollMutation = useUnenrollStudent();

    function handleSaveGrade() {
        updateGrade.mutate(
            {
                studentId: enrollment.studentId,
                courseCode: enrollment.courseCode,
                grade: gradeInput,
            },
            {
                onSuccess: () =>
                    toast.success(
                        `Updated grade for ${enrollment.course?.title ?? "course"
                        }.`
                    ),
                onError: (err) =>
                    toast.error(apiErrorMessage(err)),
            }
        );
    }

    function handleRemove() {
        if (
            !window.confirm(
                `Remove ${enrollment.course?.title ?? "this course"
                } from ${enrollment.student?.name}?`
            )
        ) {
            return;
        }

        unenrollMutation.mutate(
            {
                studentId: enrollment.studentId,
                courseCode: enrollment.courseCode,
            },
            {
                onSuccess: () =>
                    toast.success(
                        `Removed ${enrollment.course?.title ?? "course"
                        }.`
                    ),
                onError: (err) =>
                    toast.error(apiErrorMessage(err)),
            }
        );
    }

    return (
        <TableRow>

            <TableCell>

                <Badge
                    variant="outline"
                    className="font-mono"
                >
                    {enrollment.courseCode}
                </Badge>

            </TableCell>

            <TableCell className="font-medium">
                {enrollment.course?.title ?? "—"}
            </TableCell>

            <TableCell>
                {enrollment.course?.professor?.name}
            </TableCell>

            <TableCell>
                {enrollment.course?.credits ?? "—"}
            </TableCell>

            <TableCell>

                <div className="flex items-center gap-2">

                    <Select
                        value={gradeInput || "NONE"}
                        onValueChange={(value) =>
                            setGradeInput(
                                value === "NONE" ? "" : value
                            )
                        }
                    >
                        <SelectTrigger className="w-20">
                            <SelectValue placeholder="-" />
                        </SelectTrigger>

                        <SelectContent>

                            <SelectItem value="NONE">
                                -
                            </SelectItem>

                            {GRADE_OPTIONS.map((grade) => (
                                <SelectItem
                                    key={grade}
                                    value={grade}
                                >
                                    {grade}
                                </SelectItem>
                            ))}

                        </SelectContent>

                    </Select>

                    <Button
                        variant="outline"
                        size="sm"
                        onClick={handleSaveGrade}
                        disabled={
                            updateGrade.isPending ||
                            gradeInput ===
                            (enrollment.grade ?? "")
                        }
                    >
                        <Check />
                    </Button>

                </div>

            </TableCell>

            <TableCell>

                <div className="flex justify-end">

                    <Button
                        variant="destructive"
                        size="sm"
                        onClick={handleRemove}
                        disabled={
                            unenrollMutation.isPending
                        }
                    >
                        <Trash2 />

                    </Button>

                </div>

            </TableCell>

        </TableRow>
    );
}