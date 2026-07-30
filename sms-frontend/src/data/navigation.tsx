import type { ReactNode } from 'react';
import { GraduationCap, BookOpen, UserCheck2 } from 'lucide-react';

export interface NavigationItem {
  label: string;
  path?: string;
  icon?: ReactNode;
  children?: NavigationItem[];
}

export const navigationItems: NavigationItem[] = [
  {
    label: 'Roster',
    path: '/roster',
    icon: <GraduationCap size={20} />
  },
  {
    label: 'Courses',
    path: '/courses',
    icon: <BookOpen size={20} />
  },
  {
    label: 'Professors',
    path: '/professors',
    icon: <UserCheck2 size={20}/>
  }
];
