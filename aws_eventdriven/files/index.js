exports.handler = async (event) => {
  console.log(event);

  const message = event.message || process.env.defaultMessage;

  console.log("The value of message is:", message);

  return {
      statusCode: 200,
      body: JSON.stringify({
          message: message,
          success: true,
      })
  };
};
