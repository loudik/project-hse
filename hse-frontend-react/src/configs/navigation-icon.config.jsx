import {
    PiHouseLineDuotone,
    PiArrowsInDuotone,
    PiBookOpenUserDuotone,
    PiBookBookmarkDuotone,
    PiAcornDuotone,
    PiBagSimpleDuotone,
    PiGaugeDuotone,
    PiShieldCheckDuotone,
    PiFilePlusDuotone,
    PiClipboardTextDuotone,
    PiGearDuotone,
    PiUsersDuotone,
    PiKeyDuotone,
} from 'react-icons/pi'

const navigationIcon = {
    home: <PiHouseLineDuotone />,
    singleMenu: <PiAcornDuotone />,
    collapseMenu: <PiArrowsInDuotone />,
    groupSingleMenu: <PiBookOpenUserDuotone />,
    groupCollapseMenu: <PiBookBookmarkDuotone />,
    groupMenu: <PiBagSimpleDuotone />,
    // Icons for HSE menu items (from the database)
    dashboard: <PiGaugeDuotone />,
    shieldCheck: <PiShieldCheckDuotone />,
    filePlus: <PiFilePlusDuotone />,
    clipboardList: <PiClipboardTextDuotone />,
    settings: <PiGearDuotone />,
    users: <PiUsersDuotone />,
    keyRound: <PiKeyDuotone />,
}

export default navigationIcon
