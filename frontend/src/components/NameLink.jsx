import React from "react";

import Link from '@mui/material/Link';

// Shared style for the "Name" column links across the Groups/Projects
// tables: a bit larger than the surrounding body text, and using the
// theme's text color instead of the browser's default (very visible,
// hard-to-read) visited-link purple.
function NameLink({ href, children }) {
  return (
    <Link
      href={href}
      underline="hover"
      sx={{
        fontSize: "1.05rem",
        color: "text.primary",
        "&:visited": {
          color: "text.primary",
        },
      }}
    >
      {children}
    </Link>
  );
}

export default NameLink;
