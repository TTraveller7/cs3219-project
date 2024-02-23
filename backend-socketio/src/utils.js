const isPending = (status) => {
  if (status === "pending_easy" 
  || status === "pending_medium"
  || status === "pending_hard") {
      return true;
  } else {
    return false;
  }
}

const notMe = (name, nameA, nameB) => {
  if (name === nameA) {
    return nameB;
  } else {
    return nameA;
  }
}

module.exports = {
  isPending,
  notMe
}