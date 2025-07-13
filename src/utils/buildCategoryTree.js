export const buildCategoryTree = (flatList) => {
  const idMap = {};
  const root = [];

  // 1. Создай карту id → категория
  flatList.forEach((cat) => {
    idMap[cat.id] = { ...cat, elements: [], sub_elements: [] };
  });

  // 2. Собери дерево
  flatList.forEach((cat) => {
    const parent = idMap[cat.parent_id];
    if (!parent) {
      root.push(idMap[cat.id]); // верхний уровень
    } else if (parent.parent_id === null) {
      parent.elements.push(idMap[cat.id]); // 2 уровень
    } else {
      parent.sub_elements.push(idMap[cat.id]); // 3 уровень
    }
  });

  return root;
};
