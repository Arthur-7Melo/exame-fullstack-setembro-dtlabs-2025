export const getErrorMessage = (error: any): string => {
  let message = 'Server Error. Please try again.';

  if (error.response?.data?.message) {
    message = error.response.data.message;
  } else if (error.response?.data?.details) {
    message = error.response.data.details;
  } else if (error.message) {
    message = error.message;
  }

  const isShortMessage = message.length < 50;
  const isErrorCode = /^[A-Z0-9_]+$/.test(message);

  return isShortMessage || isErrorCode ? message.toUpperCase() : message;
};