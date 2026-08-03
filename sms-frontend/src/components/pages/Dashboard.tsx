// dashboard page should have cards with summary statistics, charts, and other relevant information for the user
// total students, average gpa, total courses offered, total professors, etc.
import PageHeader from "../ui/PageHeader";
import { Card } from "../ui/Card";
import { UserCheck2, Users, BookOpen, Award, GraduationCap, Star } from "lucide-react";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "../ui/Table";
import { useStudents } from "../../hooks/useStudents";
import { useCourses } from "../../hooks/useCourses";
import { useProfessors } from "../../hooks/useProfessors";
import { useStudentEnrollments } from "../../hooks/useEnrollments";

const rainbowAccents = [
    { bg: '#FEE2E2', text: '#B91C1C' },
    { bg: '#FFEDD5', text: '#C2410C' },
    { bg: '#FEF3C7', text: '#B45309' },
    { bg: '#DCFCE7', text: '#166534' },
    { bg: '#DBEAFE', text: '#1D4ED8' },
    { bg: '#EDE9FE', text: '#6D28D9' },
];

export function Dashboard() {
    const { data: students = [] } = useStudents();
    const { data: courses = [] } = useCourses();
    const { data: professors = [] } = useProfessors();
    const { data: enrollments = [] } = useStudentEnrollments(null);

    const totalStudents = students.length;
    const totalCourses = courses.length;
    const totalProfessors = professors.length;
    const activeProfessors = new Set(
        courses
            .filter((course) => Boolean(course.professorId))
            .map((course) => course.professorId as string)
    ).size;
    const averageGpa = totalStudents
        ? (students.reduce((sum, student) => sum + (student.gpa ?? 0), 0) / totalStudents).toFixed(2)
        : '0.00';
    
        //replace with revenue when financial data is available
    const honorRollStudents = students.filter((student) => (student.gpa ?? 0) >= 3.5).length;

    const recentActions = [
        ...students.map((student) => ({
            action: 'Student Added',
            details: `${student.name} was added to the roster`,
            date: student.createdAt ?? new Date().toISOString(),
        })),
        ...courses.map((course) => ({
            action: 'Course Added',
            details: `${course.title} (${course.code}) was added to the catalog`,
            date: course.createdAt ?? new Date().toISOString(),
        })),
        ...professors.map((professor) => ({
            action: 'Professor Added',
            details: `${professor.name} (${professor.id}) joined the faculty`,
            date: professor.createdAt ?? new Date().toISOString(),
        })),
        ...enrollments.map((enrollment) => ({
            action: 'Enrollment',
            details: `${enrollment.student?.name ?? `Student #${enrollment.studentId}`} enrolled in ${enrollment.course?.code ?? enrollment.courseCode}`,
            date: enrollment.enrolledAt ?? enrollment.updatedAt ?? new Date().toISOString(),
        })),
    ]
        .sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())
        .slice(0, 5);

    const stats = [
        {
            label: 'Total Students',
            value: totalStudents.toString(),
            icon: Users,
            accent: rainbowAccents[0],
        },
        {
            label: 'Total Courses',
            value: totalCourses.toString(),
            icon: BookOpen,
            accent: rainbowAccents[1],
        },
        {
            label: 'Total Professors',
            value: totalProfessors.toString(),
            icon: UserCheck2,
            accent: rainbowAccents[2],
        },
        {
            label: 'Average GPA',
            value: averageGpa,
            icon: Award,
            accent: rainbowAccents[3],
        },
        {
            label: 'Active Professors',
            value: activeProfessors.toString(),
            icon: GraduationCap,
            accent: rainbowAccents[4],
        },
        {
            label: 'Honor Roll',
            value: honorRollStudents.toString(),
            icon: Star,
            accent: rainbowAccents[5],
        },
    ];

    return (
        <>
            <div>
                <PageHeader
                    title="Dashboard"
                    description="Example University"
                />
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
                {stats.map((stat) => {
                    const Icon = stat.icon;
                    const accent = stat.accent;

                    return (
                        <Card key={stat.label} padding="lg" className="space-y-4">
                            <div
                                className="w-12 h-12 rounded-full flex items-center justify-center font-bold"
                                style={{ backgroundColor: accent.bg, color: accent.text }}
                            >
                                <Icon size={20} />
                            </div>
                            <div className="space-y-1">
                                <h2 className="text-2xl font-semibold leading-tight">{stat.value}</h2>
                                <h3 className="text-sm text-muted-foreground">{stat.label}</h3>
                            </div>
                        </Card>
                    );
                })}
            </div>
            <div>
                <Card padding="responsive">
                    <div className="mb-6">
                        <h2 className="text-xl font-semibold">
                            Recent Actions
                        </h2>
                    </div>
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead>Action</TableHead>
                                <TableHead>Details</TableHead>
                                <TableHead>Date</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {recentActions.map((item) => (
                                <TableRow key={`${item.action}-${item.details}-${item.date}`}>
                                    <TableCell>{item.action}</TableCell>
                                    <TableCell>{item.details}</TableCell>
                                    <TableCell>{new Date(item.date).toLocaleDateString()}</TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </Card>
            </div>
        </>
    )
}