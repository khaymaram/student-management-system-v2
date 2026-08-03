import PageHeader from "../ui/PageHeader";
import { Card } from "../ui/Card";
import { CircleDollarSign, Star, AlertCircle } from "lucide-react";

const rainbowAccents = [
    { bg: '#FEE2E2', text: '#B91C1C' },
    { bg: '#FFEDD5', text: '#C2410C' },
    { bg: '#FEF3C7', text: '#B45309' },
    { bg: '#DCFCE7', text: '#166534' },
    { bg: '#DBEAFE', text: '#1D4ED8' },
    { bg: '#EDE9FE', text: '#6D28D9' },
];

const stats = [
    {
        label: 'Total Collected Fees',
        icon: CircleDollarSign,
        accent: rainbowAccents[3],
    },
    {
        label: 'Outstanding Payments',
        icon: AlertCircle,
        accent: rainbowAccents[0],
    },
    {
        label: 'Total Scholarships',
        icon: Star,
        accent: rainbowAccents[2],
    }
];

export function FinancesView() {
    return (
        <>
            <div>
                <PageHeader
                    title="Finances"
                    description="View and manage school fees, payments, scholarships, and other financial records."
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
                                {/* <h2 className="text-2xl font-semibold leading-tight">{stat.value}</h2> */}
                                <h3 className="text-sm text-muted-foreground">{stat.label}</h3>
                            </div>
                        </Card>
                    );
                })}
            </div>
        </>
    );
}