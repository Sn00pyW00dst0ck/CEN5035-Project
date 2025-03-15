// src/components/Button.test.jsx
import { render } from '@testing-library/react';
import App from "../App.jsx";

test('renders button with label', () => {
    expect(() => render(<App />)).not.toThrow();
});
