// dashboard page should have cards with summary statistics, charts, and other relevant information for the user
// total students, average gpa, total courses offered, total professors, etc.
import PageHeader from "../ui/PageHeader";
import { Card } from "../ui/Card";
import { UserCheck2, Users, BookOpen, CircleDollarSign, Award } from "lucide-react";
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from "../ui/Table";

export function Dashboard() {
    return (
        <>
            <div>
                <PageHeader
                    title="Dashboard"
                    description="Example University"
                />
            </div>
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 mb-4">
                <Card padding="lg">
                    <div className="w-10 h-10 bg-sidebar-accent rounded-full flex items-center justify-center text-primary font-bold">
                        <Users  />
                    </div>
                    <h2 className="text-lg font-semibold"># ---</h2>
                    <h3 className="text-sm text-muted-foreground">Total Students</h3>
                </Card>
                <Card padding="lg">
                   <div className="w-10 h-10 bg-sidebar-accent rounded-full flex items-center justify-center text-primary font-bold">
                        <BookOpen  />
                    </div>
                    <h2 className="text-lg font-semibold"># ---</h2>
                    <h3 className="text-sm text-muted-foreground">Total Courses</h3>
                </Card>
                <Card padding="lg">
                    <div className="w-10 h-10 bg-sidebar-accent rounded-full flex items-center justify-center text-primary font-bold">
                        <UserCheck2  />
                    </div>
                    <h2 className="text-lg font-semibold"># ---</h2>
                    <h3 className="text-sm text-muted-foreground">Total Professors</h3>
                </Card>
                <Card padding="lg">
                    <div className="w-10 h-10 bg-sidebar-accent rounded-full flex items-center justify-center text-primary font-bold">
                        <CircleDollarSign  />
                    </div>
                    <h2 className="text-lg font-semibold">$ ---</h2>
                    <h3 className="text-sm text-muted-foreground">Monthly Revenue</h3>
                </Card>
                <Card padding="lg">
                    <div className="w-10 h-10 bg-sidebar-accent rounded-full flex items-center justify-center text-primary font-bold">
                        <UserCheck2  />
                    </div>
                    <h2 className="text-lg font-semibold">Active Professors %</h2>
                    <h3 className="text-sm text-muted-foreground">Active Professors</h3>
                </Card>
                <Card padding="lg">
                    <div className="w-10 h-10 bg-sidebar-accent rounded-full flex items-center justify-center text-primary font-bold">
                        <Award  />
                    </div>
                    <h2 className="text-lg font-semibold">---</h2>
                    <h3 className="text-sm text-muted-foreground">Average GPA</h3>
                </Card>
            </div>
            <div>
                {/* recent actions table: recent enrollments, recent course additions, recent professor additions, etc. */}
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
                            {/* Example rows - replace with actual data */}
                            <TableRow>
                                <TableCell>Enrollment</TableCell>
                                <TableCell>John Doe enrolled in CS101</TableCell>
                                <TableCell>2023-10-01</TableCell>
                            </TableRow>
                            <TableRow>
                                <TableCell>Course Added</TableCell>
                                <TableCell>Math 101 added to curriculum</TableCell>
                                <TableCell>2023-10-02</TableCell>
                            </TableRow>
                        </TableBody>
                    </Table>
                </Card>

            </div>
        </>
    )
}