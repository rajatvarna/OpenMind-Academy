import { render, screen } from '@testing-library/react';
import CourseCard from '../CourseCard';

describe('CourseCard', () => {
  it('renders the course title and description', () => {
    const course = {
      id: '1',
      title: 'Test Course',
      description: 'This is a test course.',
    };
    render(<CourseCard course={course} />);

    const title = screen.getByRole('heading', { name: /Test Course/i });
    expect(title).toBeInTheDocument();

    const description = screen.getByText(/This is a test course./i);
    expect(description).toBeInTheDocument();
  });
});
