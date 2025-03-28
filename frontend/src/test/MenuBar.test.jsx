import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import MenuBar from '../../src/MainScreen/ServerAndMembers/ActiveServer/MenuBar/MenuBar.jsx';
import * as React from 'react';

vi.mock('../../src/MainScreen/ServerList/ServerBadge/ServerBadge.jsx', () => ({
  default: ({ server }) => <div data-testid="server-badge">{server.name}</div>
}));

vi.mock('../../src/CommonComponents/Search/Search.jsx', () => ({
  default: ({ label }) => <div data-testid="search-component">{label}</div>
}));

vi.mock('@mui/material', () => ({
  Button: ({ children, onClick, sx }) => {
    
    if (children && children.type && children.type.name === 'default') {
      return (
        <button data-testid="members-toggle-button" onClick={onClick}>
          {children}
        </button>
      );
    }
    return (
      <button 
        data-testid={typeof children === 'string' ? "channel-button" : "members-toggle-button"} 
        onClick={onClick}
      >
        {children}
      </button>
    );
  },
  Paper: ({ children, sx }) => <div data-testid="paper-component">{children}</div>,
  Menu: ({ children, anchorEl, open, onClose }) => 
    open ? <div data-testid="menu-component">{children}</div> : null,
  MenuItem: ({ children, onClick }) => 
    <div data-testid="menu-item" onClick={onClick}>{children}</div>
}));

vi.mock('@mui/icons-material/Menu', () => ({
  default: () => <div data-testid="menu-icon">MenuIcon</div>
}));

describe('MenuBar Component', () => {
  let setVisibleMock;
  
  beforeEach(() => {
    setVisibleMock = vi.fn();
  });
  
  it('renders without crashing', () => {
    render(<MenuBar setVisible={setVisibleMock} selectedServer={{ name: 'Alice', channels: [] }} />);
    expect(screen.getByTestId('paper-component')).toBeDefined();
  });
  
  it('renders ServerBadge with correct props', () => {
    render(<MenuBar setVisible={setVisibleMock} selectedServer={{ name: 'Alice', channels: [] }} />);
    const serverBadge = screen.getByTestId('server-badge');
    expect(serverBadge).toBeDefined();
    expect(serverBadge.textContent).toBe('Alice');
  });
  
  it('renders Search component with correct label', () => {
    render(<MenuBar setVisible={setVisibleMock} selectedServer={{ name: 'Alice', channels: [] }} />);
    const searchComponent = screen.getByTestId('search-component');
    expect(searchComponent).toBeDefined();
    expect(searchComponent.textContent).toBe('Search messages');
  });
  
  it('renders menu button with icon', () => {
    render(<MenuBar setVisible={setVisibleMock} selectedServer={{ name: 'Alice', channels: [] }} />);
    const menuButton = screen.getByTestId('members-toggle-button');
    expect(menuButton).toBeDefined();
    const menuIcon = screen.getByTestId('menu-icon');
    expect(menuIcon).toBeDefined();
  });
  
  it('calls setVisible with opposite value when menu button is clicked', () => {
    render(<MenuBar setVisible={setVisibleMock} selectedServer={{ name: 'Alice', channels: [] }} />);
    const menuButton = screen.getByTestId('members-toggle-button');
    
    
    fireEvent.click(menuButton);
    expect(setVisibleMock).toHaveBeenCalledTimes(1);
    expect(setVisibleMock).toHaveBeenCalledWith(expect.any(Function));
    
    const updateFn = setVisibleMock.mock.calls[0][0];
    expect(updateFn(true)).toBe(false);
    expect(updateFn(false)).toBe(true);
  });
});