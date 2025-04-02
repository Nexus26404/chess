export enum RoomStatus {
    Unready,
    Ready,
    Started,
    Finished
}

export enum PlayerTurn {
    None,
    Black,
    White
}

export type UserInfo = {
    id: string;
    username: string;
    nickname: string;
    gamesPlayed: number;
    wins: number;
    losses: number;
    draws: number;
}
