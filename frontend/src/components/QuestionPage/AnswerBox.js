import React from "react";
import Editor from "react-simple-code-editor";
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import "prismjs/components/prism-java";
import "prismjs/components/prism-python";
import "prismjs/themes/prism.css";

import "./answerBoxStyle.css";
import { saveDocument, sendChanges } from "../../utils/SocketClientIo";

const highlightWithLineNumbers = (input, language) =>
	highlight(input, language)
		.split("\n")
		.map((line, i) => `<span class='editorLineNumber'>${i + 1}</span>${line}`)
		.join("\n");

function AnswerBox({ answer, setAnswer, socketCollabServiceClient }) {
	function changeAnswerHandler(answer) {
		setAnswer(answer);
		socketCollabServiceClient.emit(sendChanges, answer)
		socketCollabServiceClient.emit(saveDocument, answer);
	}

	return (
		<Editor
			value={answer}
			onValueChange={changeAnswerHandler}
			highlight={answer => highlightWithLineNumbers(answer, languages.js)}
			padding={10}
			textareaId="codeArea"
			className="editor"
			style={{
				fontFamily: '"Fira code", "Fira Mono", monospace',
				fontSize: 18,
				outline: 0,
				minHeight: '90vh',
			}}
		/>
	);
}

export default AnswerBox
