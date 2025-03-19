import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import MenuBar from '../../src/MainScreen/ServerAndMembers/ActiveServer/MenuBar/MenuBar.jsx';
import * as React from 'react';

// Mock the imported components
vi.mock('../../src/MainScreen/ServerList/ServerBadge/ServerBadge.jsx', () => ({
  default: ({ server }) => <div data-testid="server-badge">{server.name}</div>
}));

vi.mock('../../src/CommonComponents/Search/Search.jsx', () => ({
  default: ({ label }) => <div data-testid="search-component">{label}</div>
}));

vi.mock('@mui/material', () => ({
  Button: ({ children, onClick, sx }) => (
    <button data-testid="menu-button" onClick={onClick}>
      {children}
    </button>
  ),
  Paper: ({ children, sx }) => <div data-testid="paper-component">{children}</div>
}));

vi.mock('@mui/icons-material/Menu', () => ({
  default: () => <div data-testid="menu-icon">MenuIcon</div>
}));

describe('MenuBar Component', () => {
  let setVisibleMock;
  
  beforeEach(() => {
    // Create a mock function for setVisible
    setVisibleMock = vi.fn();
  });

  it('renders without crashing', () => {
    render(<MenuBar setVisible={setVisibleMock} />);
    expect(screen.getByTestId('paper-component')).toBeDefined();
  });

  it('renders ServerBadge with correct props', () => {
    render(<MenuBar setVisible={setVisibleMock} />);
    const serverBadge = screen.getByTestId('server-badge');
    expect(serverBadge).toBeDefined();
    expect(serverBadge.textContent).toBe('Alice');
  });

  it('renders Search component with correct label', () => {
    render(<MenuBar setVisible={setVisibleMock} />);
    const searchComponent = screen.getByTestId('search-component');
    expect(searchComponent).toBeDefined();
    expect(searchComponent.textContent).toBe('Search messages');
  });

  it('renders menu button with icon', () => {
    render(<MenuBar setVisible={setVisibleMock} />);
    const menuButton = screen.getByTestId('menu-button');
    expect(menuButton).toBeDefined();
    const menuIcon = screen.getByTestId('menu-icon');
    expect(menuIcon).toBeDefined();
  });

  it('calls setVisible with opposite value when menu button is clicked', () => {
    render(<MenuBar setVisible={setVisibleMock} />);
    const menuButton = screen.getByTestId('menu-button');
    
    // First click
    fireEvent.click(menuButton);
    expect(setVisibleMock).toHaveBeenCalledTimes(1);
    expect(setVisibleMock).toHaveBeenCalledWith(expect.any(Function));
    
    // Check that the function passed to setVisible returns the opposite of its input
    const updateFn = setVisibleMock.mock.calls[0][0];
    expect(updateFn(true)).toBe(false);
    expect(updateFn(false)).toBe(true);
  });
});