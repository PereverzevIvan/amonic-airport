import { useQuery } from "@tanstack/react-query";
import axios from "axios";
import { useEffect, useState } from "react";
import "./SurvayPage.scss";
import Button from "../../components/Button/Button";

const BEGIN_DATE = "2016-01-01";
// const BEGIN_DATE = "2025-01-01";
const END_DATE = "2025-01-01";

const RESULT_DIAGRAM_COLORS = [
  "#0aff00",
  "#9eff9d",
  "#008c53",
  "#ffff00",
  "#ffa500",
  "#ff0000",
  "#808080",
];

const getGroupsWithValues = async () => {
  var response = await axios.get(
    `http://localhost:3000/api/survey/groups-with-values`
  );

  // Delete Departure
  Object.entries(response.data).forEach(([key, group]) => {
    if (group.name === "Departure") {
      delete response.data[key];
      return;
    }
  });
  return response.data;
};

const getQuestionsWithAnswers = async () => {
  var response = await axios.get(
    `http://localhost:3000/api/survey/questions-with-answers`
  );
  return response.data;
};

const initDefaultRespondentsAnswers = (new_data_questions, new_data_groups) => {
  var data_respondents_answers = {};
  // questions
  Object.entries(new_data_questions).map(([question_id, question]) => {
    data_respondents_answers[question_id] = {};

    // answers
    question.answers.map((answer, _) => {
      data_respondents_answers[question_id][answer.id] = {};
      // data_respondents_answers[question_name][idx_answer] = {};

      // groups
      Object.entries(new_data_groups).map(([group_id, group]) => {
        const group_values = group.values;
        data_respondents_answers[question_id][answer.id][group_id] = {};

        // group_values
        group_values.map((group_value) => {
          // count_answers for group_value
          // default = 0
          data_respondents_answers[question_id][answer.id][group_id][
            group_value.id
          ] = 0;
        });
      });
    });
  });

  return data_respondents_answers;
};

const generateFilterGroupValues = (filter_group_values) => {
  var query_filter_group_values = "";
  if (filter_group_values) {
    Object.entries(filter_group_values).map(([_, filter]) => {
      if (!filter.is_active || filter.group_value_id == 0) return;

      query_filter_group_values += `&filter_group_values=${filter.group_value_id}`;
    });
  }

  console.log("query_filter_group_values:", query_filter_group_values);
  return query_filter_group_values;
};

const getRespondentsAnswers = async ({ queryKey }) => {
  const [
    _,
    new_data_questions,
    new_data_groups,
    begin_date,
    end_date,
    filter_group_values,
  ] = queryKey;

  // console.log("data_respondents_answers:", new_data_questions, new_data_groups);

  try {
    var data_respondents_answers = initDefaultRespondentsAnswers(
      new_data_questions,
      new_data_groups
    );

    var response = await axios.get(
      `http://localhost:3000/api/surveys/respondents-answers?` +
        `begin_date=${begin_date}` +
        `&end_date=${end_date}` +
        generateFilterGroupValues(filter_group_values)
    );

    response.data.map((respondent_answer) => {
      var question_id = respondent_answer.question_id;
      var answer_id = respondent_answer.answer_id;
      var group_id = respondent_answer.group_id;
      var group_value_id = respondent_answer.group_value_id;
      var count_answers = respondent_answer.count_answers;

      // question_id: 1
      // answer_id: 115
      // group_id: 1
      // group_value_id: 39
      // count_answers: 1

      // console.log(question_id, answer_id, group_id, group_value_id, count_answers);

      // discard invalid groups by checking if group_id in data_respondents_answers[question_id][answer_id
      if (group_id in data_respondents_answers[question_id][answer_id])
        data_respondents_answers[question_id][answer_id][group_id][
          group_value_id
        ] += count_answers;
    });
  } catch (error) {
    console.error(error);
  }

  console.log(response.data);

  return data_respondents_answers;
};

const initFilterGroupValues = (data_groups, setFilterGroupValues) => {
  var filter_group_values = {};
  if (!data_groups) {
    return;
  }

  Object.entries(data_groups).map(([group_id]) => {
    filter_group_values[group_id] = {
      is_active: true,
      group_value_id: 0,
    };
  });

  setFilterGroupValues(filter_group_values);
};

const updateCountActiveGroupValueColumns = (
  data_groups,
  filter_group_values,
  setCountActiveGroupValueColumns
) => {
  var count_active_group_value_columns = 0;

  // console.log("aboba", data_groups, filter_group_values, filter_group_values.length);

  Object.entries(data_groups).map(([group_id, group]) => {
    if (filter_group_values[group_id].is_active == false) return;
    if (filter_group_values[group_id].group_value_id == 0) {
      count_active_group_value_columns += group.values.length;
    } else {
      count_active_group_value_columns += 1;
    }
  });

  setCountActiveGroupValueColumns(count_active_group_value_columns);
};

function SurveyPage() {
  const [showFull, setShowFull] = useState(false);
  var [begin_date, setBeginDate] = useState(BEGIN_DATE);
  var [end_date, setEndDate] = useState(END_DATE);
  var [filter_group_values, setFilterGroupValues] = useState({});
  var [count_active_group_value_columns, setCountActiveGroupValueColumns] =
    useState(0);

  const {
    data: data_groups,
    error: errorGroups,
    isLoading: isLoadingGroups,
  } = useQuery({
    queryKey: ["groups"],
    queryFn: getGroupsWithValues,
  });

  useEffect(() => {
    if (!data_groups) return;
    initFilterGroupValues(data_groups, setFilterGroupValues);
  }, [data_groups]);

  useEffect(() => {
    if (
      !data_groups ||
      !filter_group_values ||
      Object.keys(filter_group_values).length == 0
    ) {
      return;
    }
    updateCountActiveGroupValueColumns(
      data_groups,
      filter_group_values,
      setCountActiveGroupValueColumns
    );
  }, [data_groups, filter_group_values]);

  const {
    data: data_questions,
    error: errorQuestions,
    isLoading: isLoadingQuestions,
  } = useQuery({
    queryKey: ["questions"],
    queryFn: getQuestionsWithAnswers,
  });

  const {
    data: data_respondents_answers,
    error: errorRespondentsAnswers,
    isLoading: isLoadingRespondentsAnswers,
  } = useQuery({
    queryKey: [
      "respondents_answers",
      data_questions,
      data_groups,
      begin_date,
      end_date,
      filter_group_values,
    ],
    queryFn: getRespondentsAnswers,
    enabled: !isLoadingGroups && !isLoadingQuestions, // Only run the query if userId is truthy
  });

  if (isLoadingGroups || isLoadingQuestions) return <div>Loading...</div>;
  if (errorGroups) return <div>Error: {errorGroups.message}</div>;
  if (errorQuestions) return <div>Error: {errorQuestions.message}</div>;

  if (isLoadingRespondentsAnswers)
    return <div>Loading respondents answers...</div>;
  if (errorRespondentsAnswers)
    return <div>Error: {errorRespondentsAnswers.message}</div>;
  console.log("groups:", data_groups);
  console.log("questions:", data_questions);

  // var data_groups = {
  //     "group1": ["value1", "value2"],
  //     "group2": ["value1", "value2", "value3"],
  // };

  // var data_questions = {
  //     "question1": ["ans1", "ans2", "ans3"],
  //     "question2": ["ans1", "ans2", "ans3"],
  // }

  /*
    {
        // questions
        "Q1": [ 
            // answers
            0: {
                // groups
                "group1": {
                    // group_values with count answers
                    "value1": 0,
                    "value2": 0
                }
            },
        ]
    }
    */

  // var data_respondents_answers = initDefaultRespondentsAnswers(new_data_questions, new_data_groups);
  // console.log("data_respondents_answers", data_respondents_answers);

  // setFilterGroupValues({});
  // initFilterGroupValues(data_groups, setFilterGroupValues);
  console.log("filter_group_values", filter_group_values);

  var map_respondents_answers = (question_id, answer_id, answer_idx) => {
    // console.log("grv", data_respondents_answers[question_name][answer_id]);
    // console.log("entr", Object.entries(data_respondents_answers[question_name][answer_id]));

    return Object.entries(data_respondents_answers[question_id][answer_id]).map(
      ([group_id, group_values]) => {
        if (filter_group_values[group_id].is_active == false) return;
        // console.log("check", group_values);
        return Object.entries(group_values).map(
          ([group_value_id, count_answers], group_value_idx) => {
            if (
              filter_group_values[group_id].group_value_id != 0 &&
              filter_group_values[group_id].group_value_id != group_value_id
            ) {
              return;
            }

            // const edge_border_style_str = "1px solid " + RESULT_DIAGRAM_COLORS[answer_idx]
            const edge_border_style_str = "3px solid black";

            var edge_border_styles = {};
            if (filter_group_values[group_id].group_value_id != 0) {
              edge_border_styles["borderLeft"] = edge_border_style_str;
              edge_border_styles["borderRight"] = edge_border_style_str;
            }
            if (group_value_idx == 0) {
              edge_border_styles["borderLeft"] = edge_border_style_str;
            }
            if (group_value_idx == group_values.length - 1) {
              edge_border_styles["borderRight"] = edge_border_style_str;
            }

            return (
              <td
                key={question_id + " " + answer_id + " " + group_value_id}
                style={{
                  //borderColor: RESULT_DIAGRAM_COLORS[answer_idx] + "55",
                  backgroundColor: RESULT_DIAGRAM_COLORS[answer_idx] + "33",
                  ...edge_border_styles,
                }}
              >
                {count_answers}
              </td>
            );
          }
        );
      }
    );
  };

  var map_general_respondents_answers = () => {
    return Object.entries(data_groups).map(([group_id, group]) => {
      if (filter_group_values[group_id].is_active == false) return <></>;

      return group.values.map((group_value, group_value_index) => {
        if (
          filter_group_values[group_id].group_value_id != 0 &&
          filter_group_values[group_id].group_value_id != group_value.id
        ) {
          return <></>;
        }

        if (filter_group_values[group_id].is_active == false) return <></>;

        let countGroupValueAnswers = 0;
        Object.entries(data_respondents_answers).map(([_, question]) => {
          Object.entries(question).map(([_, answer]) => {
            Object.entries(answer).map(([_, group]) => {
              if (group_value.id in group) {
                countGroupValueAnswers += group[group_value.id];
              }
            });
          });
        });

        return (
          <td
            scope="col"
            style={
              group_value_index == 0
                ? { borderLeft: "3px solid black" }
                : group_value_index == group.values.length
                ? { borderRight: "3px solid black" }
                : {}
            }
          >
            {countGroupValueAnswers}
          </td>
        );
      });
    });
  };

  var calc_row_total = (question_id, answer_id) => {
    var total_sum = 0;

    Object.entries(data_respondents_answers[question_id][answer_id]).map(
      ([group_id, group_values]) => {
        if (filter_group_values[group_id].is_active == false) return;

        Object.entries(group_values).map(([group_value_id, count_answers]) => {
          if (
            filter_group_values[group_id].group_value_id != 0 &&
            filter_group_values[group_id].group_value_id != group_value_id
          ) {
            return;
          }

          total_sum += count_answers;
        });
      }
    );

    return total_sum;
  };

  var generateResultDiagram = (question) => {
    var total_question_answers_count = 0;
    var question_answers_count_arr = [];

    question.answers.map((answer, _) => {
      var answer_row_total = calc_row_total(question.id, answer.id);
      question_answers_count_arr.push(answer_row_total);
      total_question_answers_count += answer_row_total;
    });

    return question.answers.map((_, answer_idx) => {
      var cur_background_color = RESULT_DIAGRAM_COLORS[answer_idx];

      var width_ratio =
        question_answers_count_arr[answer_idx] / total_question_answers_count;
      var width = width_ratio * 100 + "%";

      return (
        <div
          style={{
            width: width,
            height: ".5rem",
            backgroundColor: cur_background_color,
          }}
        ></div>
      );
    });
  };

  return (
    <>
      <section className="survey-page">
        {showFull == true ? (
          <Button
            style={{ justifySelf: "start" }}
            onClick={() => setShowFull(false)}
          >
            Скрыть полную версию
          </Button>
        ) : (
          <Button
            style={{ width: "20% !important" }}
            onClick={() => setShowFull(true)}
          >
            Показать полную версию
          </Button>
        )}
      </section>

      {!showFull && (
        <section className="survey-page">
          <table className="table" style={{ width: "100%" }}>
            <caption>Результаты опроса</caption>
            <thead className="table__head">
              <tr className="table__row" key={0}>
                {/* LEFT SIDE FOR QUESTIONS AND STATS */}
                {Object.entries(data_groups).map(([group_id, group]) => {
                  if (filter_group_values[group_id].is_active == false) return;

                  return (
                    <th
                      scope="col"
                      style={{
                        borderLeft: "3px solid black",
                        borderRight: "3px solid black",
                      }}
                      colSpan={
                        filter_group_values[group_id].group_value_id == 0
                          ? group.values.length
                          : 1
                      }
                      key={"group" + group_id}
                    >
                      {group.name}
                    </th>
                  );
                })}
              </tr>
              <tr key={1}>
                {Object.entries(data_groups).map(([group_id, group]) => {
                  if (filter_group_values[group_id].is_active == false) return;

                  return group.values.map((group_value, group_value_index) => {
                    if (
                      filter_group_values[group_id].group_value_id != 0 &&
                      filter_group_values[group_id].group_value_id !=
                        group_value.id
                    ) {
                      return;
                    }
                    return (
                      <th
                        scope="col"
                        style={
                          group_value_index == 0
                            ? { borderLeft: "3px solid black" }
                            : group_value_index == group.values.length
                            ? { borderRight: "3px solid black" }
                            : {}
                        }
                      >
                        {group_value.name}
                      </th>
                    );
                  });
                })}
              </tr>
            </thead>
            <tbody key={"test-table-body"}>
              {map_general_respondents_answers()}
            </tbody>
          </table>
        </section>
      )}

      {showFull && (
        <section className="survey-page">
          <label>
            Select begin date:
            <input
              className="form__input"
              type="date"
              id="begin_date"
              value={begin_date}
              onChange={(e) => {
                console.log("change", e.target.value);
                setBeginDate(e.target.value);
                // setSelectedDate(e.target.value);
              }}
            />
          </label>
          <label>
            Select end date:
            <input
              className="form__input"
              type="date"
              id="end_date"
              value={end_date}
              onChange={(e) => {
                console.log("change", e.target.value);
                setEndDate(e.target.value);
                // setSelectedDate(e.target.value);
              }}
            />
          </label>
          <table className="table" style={{ width: "100%" }}>
            <caption>Результаты опроса</caption>
            <thead className="table__head">
              <tr className="table__row" key={0}>
                {/* LEFT SIDE FOR QUESTIONS AND STATS */}
                <th scope="col" key={0}></th>
                <th scope="col" key={1}></th>
                {Object.entries(data_groups).map(([group_id, group]) => {
                  if (filter_group_values[group_id].is_active == false) return;

                  return (
                    <th
                      scope="col"
                      style={{
                        borderLeft: "3px solid black",
                        borderRight: "3px solid black",
                      }}
                      colSpan={
                        filter_group_values[group_id].group_value_id == 0
                          ? group.values.length
                          : 1
                      }
                      key={"group" + group_id}
                    >
                      {group.name}
                    </th>
                  );
                })}
              </tr>
              <tr key={1}>
                <th scope="col"></th>
                <th scope="col">Total</th>
                {Object.entries(data_groups).map(([group_id, group]) => {
                  if (filter_group_values[group_id].is_active == false) return;

                  return group.values.map((group_value, group_value_index) => {
                    if (
                      filter_group_values[group_id].group_value_id != 0 &&
                      filter_group_values[group_id].group_value_id !=
                        group_value.id
                    ) {
                      return;
                    }
                    return (
                      <th
                        scope="col"
                        style={
                          group_value_index == 0
                            ? { borderLeft: "3px solid black" }
                            : group_value_index == group.values.length
                            ? { borderRight: "3px solid black" }
                            : {}
                        }
                      >
                        {group_value.name}
                      </th>
                    );
                  });
                })}
              </tr>
            </thead>
            <tbody key={"test-table-body"}>
              {Object.entries(data_questions).map(([_, question]) => (
                <>
                  <tr key={question.id}>
                    <th scope="row" key={question.text}>
                      {question.text}
                    </th>
                    <td
                      scope="row"
                      colSpan={1 + count_active_group_value_columns} // +1 for 'total' column
                      style={{
                        border: "none",
                        borderTop: "1px solid white",
                      }}
                    >
                      <div style={{ display: "flex", width: "100%" }}>
                        {generateResultDiagram(question)}
                      </div>
                    </td>
                  </tr>
                  {question.answers.map((answer, answer_idx) => (
                    <tr key={question.id + " " + answer.id}>
                      <td
                        scope="row"
                        key={"answer " + question.id + " " + answer.id}
                        style={{
                          textAlign: "right",
                          //borderColor: RESULT_DIAGRAM_COLORS[answer_idx] + "55",
                          backgroundColor:
                            RESULT_DIAGRAM_COLORS[answer_idx] + "33",
                        }}
                      >
                        {answer.text}
                      </td>
                      <td
                        scope="row"
                        key={question.id + " " + answer.id}
                        style={{
                          //borderColor: RESULT_DIAGRAM_COLORS[answer_idx] + "55",
                          backgroundColor:
                            RESULT_DIAGRAM_COLORS[answer_idx] + "33",
                        }}
                      >
                        <b>{calc_row_total(question.id, answer.id)}</b>
                      </td>
                      {map_respondents_answers(
                        question.id,
                        answer.id,
                        answer_idx
                      )}
                    </tr>
                  ))}
                </>
              ))}
            </tbody>
          </table>
          <div
            style={{
              display: "flex",
              width: "100%",
              gap: "1rem",
              flexWrap: "wrap",
              justifyContent: "center",
              alignItems: "center",
            }}
          >
            {data_questions[Object.keys(data_questions)[0]].answers.map(
              (answer, answer_idx) => (
                <div
                  style={{
                    display: "flex",
                    gap: ".2rem",
                    alignItems: "center",
                  }}
                >
                  <div
                    style={{
                      width: "1rem",
                      height: "1rem",
                      backgroundColor: RESULT_DIAGRAM_COLORS[answer_idx],
                    }}
                  />
                  <p>{answer.text}</p>
                </div>
              )
            )}
          </div>

          {/* Filters */}
          {Object.entries(data_groups).map(([group_id, group]) => (
            <div>
              <label>
                <input
                  className="form__input"
                  type="checkbox"
                  value={group.name}
                  checked={filter_group_values[group_id].is_active}
                  onChange={(e) => {
                    console.log("toggle filter: ", e.target.checked);
                    setFilterGroupValues({
                      ...filter_group_values,
                      [group_id]: {
                        ...filter_group_values[group_id],
                        is_active: e.target.checked,
                      },
                    });
                  }}
                />
                {group.name}
              </label>
              <select
                className="form__input"
                value={filter_group_values[group_id].group_value_id}
                onChange={(e) => {
                  console.log("change filter: ", e.target.value);

                  setFilterGroupValues({
                    ...filter_group_values,
                    [group_id]: {
                      ...filter_group_values[group_id],
                      group_value_id: e.target.value,
                    },
                  });
                }}
              >
                <option value={0}>All</option>
                {group.values.map((group_value) => (
                  <option value={group_value.id}>{group_value.name}</option>
                ))}
              </select>
            </div>
          ))}
        </section>
      )}
    </>
  );
}

export default SurveyPage;
