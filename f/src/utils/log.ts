// utils/log.ts

// Override console.log to disable logging in production
export const overrideConsoleLog = () => {
    if (process.env.NODE_ENV === 'production') {
      // Disable console logs in production by replacing them with an empty function
      console.log = () => {};
      console.warn = () => {};
      console.error = () => {};
    }
  };