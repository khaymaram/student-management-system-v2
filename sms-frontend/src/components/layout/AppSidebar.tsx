import { useState } from 'react';
import { ChevronLeft, ChevronRight, GraduationCap } from 'lucide-react';
import SidebarItem from '../ui/SidebarItem';
import { navigationItems } from '@/data/navigation';
const logo = new URL('../../assets/logo.png', import.meta.url).href;

const AppSidebar = () => {
  const [isCollapsed, setIsCollapsed] = useState(false);

  return (
    <>
      {/* Sidebar */}
      <aside
        className={`
          fixed top-0 left-0 h-screen bg-sidebar text-sidebar-foreground border-r border-sidebar-border
          transition-all duration-500 ease-in-out z-40 will-change-transform
          ${isCollapsed ? 'w-20' : 'w-64'}
        `}
      >
        <div className="flex flex-col h-full">
          {/* Header */}
          <div
            className={`
            flex items-center h-16 border-b border-sidebar-border
            transition-all duration-300
            ${isCollapsed ? 'justify-center px-2' : 'justify-between pl-4 pr-2'}
          `}
          >
            {!isCollapsed && (
              <div className="flex items-center gap-3 min-w-0">
                <img src={logo} alt="SMS logo" className="w-8 h-8 object-contain rounded-md shrink-0" />
                <span className="font-heading font-bold text-sm text-sidebar-foreground whitespace-nowrap truncate">
                  GRGI University
                </span>
              </div>
            )}

            {/* Collapse Toggle */}
            <button
              onClick={() => setIsCollapsed(!isCollapsed)}
              className="flex items-center justify-center w-10 h-10 rounded-lg hover:bg-sidebar-accent transition-all duration-200 hover:scale-110 shrink-0"
              aria-label="Toggle sidebar width"
            >
              {isCollapsed ? (
                <ChevronRight size={20} className="text-sidebar-foreground/80 transition-transform" />
              ) : (
                <ChevronLeft size={20} className="text-sidebar-foreground/80 transition-transform" />
              )}
            </button>
          </div>

          {/* Navigation */}
          <nav
            className={`flex-1 overflow-y-auto py-4 scroll-smooth custom-scrollbar ${
              isCollapsed ? 'overflow-visible' : 'overflow-y-auto'
            }`}
          >
            {navigationItems.map((item, index) => (
              <SidebarItem key={index} item={item} isCollapsed={isCollapsed} />
            ))}
          </nav>

          {/* Footer */}
          <div className="p-4 border-t border-sidebar-border transition-all duration-300">
            <div className={`flex items-center gap-2 ${isCollapsed ? 'justify-center' : ''}`}>
              <div className="w-10 h-10 bg-sidebar-accent rounded-full flex items-center justify-center text-sidebar-foreground font-semibold shrink-0">
                <GraduationCap size={18} />
              </div>
              {!isCollapsed && (
                <div className="flex-1 min-w-0 text-left">
                  <p className="text-sm font-medium text-sidebar-foreground truncate">Simple SMS</p>
                  <p className="text-xs text-sidebar-foreground/70 truncate">Student Management System</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </aside>

      {/* Spacer for main content */}
      <div
        className={`
        shrink-0 transition-all duration-500 ease-in-out
        ${isCollapsed ? 'w-20' : 'w-64'}
      `}
      />
    </>
  );
};

export default AppSidebar;
