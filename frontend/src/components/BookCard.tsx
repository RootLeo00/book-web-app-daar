import { CardMedia, Card, CardHeader, Avatar, CardContent, Typography } from '@mui/material';
import red from '@mui/material/colors/red';
import avatar from '../images/avatar_icon.png';


const BookCard = ({ cardData }: any) => {

    return ((cardData !== null || cardData !== undefined ) ? (
        <Card sx={{ maxWidth: 345 }}>
            <CardHeader
                title={cardData?.title}
                subheader={(cardData?.author === undefined || cardData?.author === "None") ? ("No Author") : (cardData?.author)} //TODO to fix python backend
                avatar={
                    <Avatar
                        src={avatar}
                        sx={{ bgcolor: red[500], width: 56, height: 56 }}
                    />
                }
            />
            <CardMedia
                component="img"
                height="500"
                image={cardData?.imageBook}
                alt="Loading card..."
            />
            <CardContent>
                <Typography variant="body2" color="text.secondary">
                    {cardData?.text}
                </Typography>
                <Typography variant="caption" color="text.secondary">
                    Crank: {cardData?.crank}
                </Typography>
                <br></br>
                <Typography variant="caption" color="text.secondary">
                    Occurrence: {cardData?.occurrence}
                </Typography>
            </CardContent>
        </Card>
    ) : (<p>loading card</p>))
}

export default BookCard;