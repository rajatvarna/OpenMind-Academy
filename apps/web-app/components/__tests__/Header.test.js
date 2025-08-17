import { render, screen } from '@testing-library/react';
import Header from '../Header';
import { AuthContext } from '../../context/AuthContext';

// Mock the SearchBar and DonationModal components
jest.mock('../SearchBar', () => () => <div data-testid="search-bar" />);
jest.mock('../DonationModal', () => () => <div data-testid="donation-modal" />);

const renderWithAuth = (ui, { providerProps, ...renderOptions }) => {
  return render(
    <AuthContext.Provider value={providerProps}>{ui}</AuthContext.Provider>,
    renderOptions
  );
};

describe('Header', () => {
  it('renders the logo', () => {
    const providerProps = { user: null, loading: false, stats: { score: 0 } };
    renderWithAuth(<Header />, { providerProps });
    expect(screen.getByText('FreeEdu')).toBeInTheDocument();
  });

  it('renders the login button when user is not logged in', () => {
    const providerProps = { user: null, loading: false, stats: { score: 0 } };
    renderWithAuth(<Header />, { providerProps });
    expect(screen.getByText('Login')).toBeInTheDocument();
  });

  it('renders the user name and logout button when user is logged in', () => {
    const user = { name: 'Test User' };
    const stats = { score: 100 };
    const providerProps = { user, stats, loading: false };
    renderWithAuth(<Header />, { providerProps });
    expect(screen.getByText(`Hello, ${user.name} (⭐ ${stats.score})`)).toBeInTheDocument();
    expect(screen.getByText('Logout')).toBeInTheDocument();
  });

  it('renders the support us button', () => {
    const providerProps = { user: null, loading: false, stats: { score: 0 } };
    renderWithAuth(<Header />, { providerProps });
    expect(screen.getByText('Support Us')).toBeInTheDocument();
  });
});
