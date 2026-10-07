DELETE rma FROM role_menu_access rma
  JOIN menus m ON m.id = rma.menu_id
  JOIN roles r ON r.id = rma.role_id
  WHERE m.name = 'Dashboard' AND r.name = 'ANP HSE';