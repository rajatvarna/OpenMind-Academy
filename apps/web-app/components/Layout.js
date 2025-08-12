import React from 'react';
import Header from './Header';

const Layout = ({ children }) => {
  return (
    <>
      <Header />
      <main>
        {children}
      </main>
      <footer style={{ textAlign: 'center', padding: '2rem 0', borderTop: '1px solid #eaeaea' }}>
        <p>Free Education For All - A Demo Project</p>
      </footer>
    </>
  );
};

export default Layout;
