// page details will have: 
//  back to students button at the top of the page 
// Name
// in the header area there will be: number of courses enrolled, grade level, credits, gpa
// body: table with courses enrolled: id, title, course grade, unenroll from course button
import { useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { ArrowLeft, Check, Trash2, Plus, HandCoins } from "lucide-react";
import { toast } from "sonner";
import { EnrollmentModal } from "./EnrollmentModal";
import Modal from "../ui/Modal";
import Input from "../ui/Input";

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

import { useFinance, useUpdateFinance } from "../../hooks/useFinances";

import { apiErrorMessage } from "../../lib/axios";

import type { Student, Enrollment, Finance } from "../../types";

const GRADE_OPTIONS = ["A", "B", "C", "D", "F"] as const;

const currency = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
    maximumFractionDigits: 0,
});

const rainbowAccents = [
    { bg: '#FEE2E2', text: '#B91C1C' },
    { bg: '#FFEDD5', text: '#C2410C' },
    { bg: '#FEF3C7', text: '#B45309' },
    { bg: '#DCFCE7', text: '#166534' },
    { bg: '#DBEAFE', text: '#1D4ED8' },
    { bg: '#EDE9FE', text: '#6D28D9' },
];

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
        data: financeData,
        isLoading: financeLoading,
    } = useFinance(parsedStudentId)

    const finance = financeData as Finance | undefined;

    const [paymentFinance, setPaymentFinance] = useState<{
        studentId: number;
        paid: number;
        scholarship: number;
        isInState: boolean;
    } | null>(null);
    const updateFinance = useUpdateFinance();
    const [paymentAmount, setPaymentAmount] = useState("");

    const openPayModal = (finance: { studentId: number; paid: number; scholarship: number; isInState: boolean }) => {
        setPaymentFinance(finance);
        setPaymentAmount("");
    };

    const closePayModal = () => {
        setPaymentFinance(null);
        setPaymentAmount("");
    };

    const handlePay = () => {
        if (!paymentFinance) return;

        const amount = Number(paymentAmount);
        if (Number.isNaN(amount) || amount <= 0) {
            toast.error("Enter a valid payment amount greater than 0.");
            return;
        }

        updateFinance.mutate(
            {
                studentId: paymentFinance.studentId,
                input: {
                    scholarship: paymentFinance.scholarship,
                    paid: paymentFinance.paid + amount,
                    isInState: paymentFinance.isInState,
                },
            },
            {
                onSuccess: () => {
                    toast.success("Payment recorded.");
                    closePayModal();
                },
                onError: (error) => toast.error(apiErrorMessage(error)),
            }
        );
    };

    const {
        data: enrollments,
        isLoading: enrollmentsLoading,
    } = useStudentEnrollments(parsedStudentId);

    const totalCredits = useMemo(
        () => getTotalCredits(enrollments),
        [enrollments]
    );

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


            {/* Finances */}
            {finance && !financeLoading && student && (
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4">
                    <Card padding={"responsive"}>
                        <div className="mb-6">
                            <h2 className="text-xl font-semibold">
                                {student.name}'s Finances
                            </h2>
                            <div>
                                <Badge style={{
                                    backgroundColor: finance.isInState ? rainbowAccents[4].bg : rainbowAccents[5].bg,
                                    color: finance.isInState ? rainbowAccents[4].text : rainbowAccents[5].text,
                                }}>
                                    {finance.isInState ? "In-state" : "Out-of-state"}
                                </Badge>
                            </div>
                        </div>

                        <Table >
                            <TableBody>
                                <TableRow >
                                    <TableCell>Tuition:</TableCell>
                                    <TableCell align="right">{currency.format(finance.tuition)}</TableCell>
                                </TableRow>
                                <TableRow>
                                    <TableCell>Scholarships:</TableCell>
                                    <TableCell align="right">{currency.format(finance.scholarship)}</TableCell>
                                </TableRow>
                                <TableRow>
                                    <TableCell>Paid:</TableCell>
                                    <TableCell align="right">{currency.format(finance.paid)}</TableCell>
                                </TableRow>
                                <TableRow style={{ backgroundColor: finance.remaining > 0 ? rainbowAccents[0].bg : rainbowAccents[3].bg, color: finance.remaining > 0 ? rainbowAccents[0].text : rainbowAccents[3].text }}>
                                    <TableCell>Balance Remaining:</TableCell>
                                    <TableCell align="right">{currency.format(finance.remaining)}</TableCell>
                                </TableRow>
                            </TableBody>
                        </Table>

                        <br></br>
                        <Button type="button" onClick={() => openPayModal(finance)}>
                            <HandCoins />
                            Pay
                        </Button>
                    </Card>
                    <Card padding={"responsive"}>
                        <div>
                            Pie Chart with Remaining, Paid, Scholarship
                        </div>
                    </Card>
                </div>

            )}

            <EnrollmentModal student={enrollmentStudent} onClose={() => setEnrollmentStudent(null)} />
            <Modal isOpen={!!paymentFinance} onClose={closePayModal} title="Pay Student Fees">
                <div className="space-y-4">
                    <Input
                        label="Payment Amount ($)"
                        inputMode="decimal"
                        value={paymentAmount}
                        onChange={(event) => setPaymentAmount(event.target.value)}
                        placeholder={finance?.remaining ? `Remaining Balance: ${currency.format(finance.remaining)}` : ""}
                    />

                    <div className="flex justify-end gap-2 pt-2">
                        <Button type="button" variant="outline" onClick={closePayModal}>
                            Cancel
                        </Button>
                        <Button type="button" onClick={handlePay} disabled={updateFinance.isPending}>
                            {updateFinance.isPending ? "Processing..." : "Confirm Payment"}
                        </Button>
                    </div>
                </div>
            </Modal>
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
                {/* the professor name is not showing up */}
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