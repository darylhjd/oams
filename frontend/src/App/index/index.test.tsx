import React from 'react';
import {Component} from './index'
import {render, screen} from '@testing-library/react';

test('renders guest index', () => {
  render(<Component/>);
  const div = screen.getByText('OAMS')
  expect(div).toBeInTheDocument();
});
